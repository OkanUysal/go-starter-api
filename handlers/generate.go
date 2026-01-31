package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/OkanUysal/go-logger"
	"github.com/OkanUysal/go-starter-api/generator"
	"github.com/OkanUysal/go-starter-api/types"
	"github.com/gin-gonic/gin"
)

// GenerateProject generates a new project and returns a ZIP file
// @Summary      Generate a new Go project
// @Description  Generates a new Go project with selected libraries and configuration, returns a ZIP file
// @Tags         Generator
// @Accept       json
// @Produce      application/zip
// @Param        request  body      types.GenerateRequest  true  "Project configuration"
// @Success      200      {file}    binary                 "ZIP file download"
// @Failure      400      {object}  types.GenerateResponse "Bad request"
// @Failure      500      {object}  types.GenerateResponse "Internal server error"
// @Router       /generate [post]
func GenerateProject(c *gin.Context) {
	var req types.GenerateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid request body", logger.Err(err))
		c.JSON(400, types.GenerateResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	logger.Info("Generating project", logger.String("name", req.Name), logger.String("modulePath", req.ModulePath))

	// Validate request
	if req.Name == "" {
		logger.Warn("Project name is required")
		c.JSON(400, types.GenerateResponse{
			Success: false,
			Error:   "Project name is required",
		})
		return
	}

	if req.ModulePath == "" {
		logger.Warn("Module path is required")
		c.JSON(400, types.GenerateResponse{
			Success: false,
			Error:   "Module path is required",
		})
		return
	}

	// Create temporary directory for project
	tempDir := filepath.Join("temp", fmt.Sprintf("%s_%d", req.Name, time.Now().Unix()))
	projectDir := filepath.Join(tempDir, req.Name)

	defer func() {
		// Cleanup temp directory after some time
		time.AfterFunc(10*time.Minute, func() {
			os.RemoveAll(tempDir)
		})
	}()

	// Generate project
	config := generator.ProjectConfig{
		Name:       req.Name,
		ModulePath: req.ModulePath,
		Structure:  req.Structure,
		Database:   req.Database.Type,
		Libraries:  req.Libraries,
		Deployment: req.Deployment,
		OutputDir:  projectDir,
	}

	if err := generator.GenerateProject(&config); err != nil {
		logger.Error("Failed to generate project", logger.Err(err), logger.String("project", req.Name))
		if Metrics != nil {
			Metrics.IncrementCounter("projects_generated_total", map[string]string{"status": "failed"})
		}
		c.JSON(500, types.GenerateResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to generate project: %v", err),
		})
		return
	}

	logger.Info("Project generated successfully", logger.String("project", req.Name))

	// Create ZIP file
	zipFileName := fmt.Sprintf("%s.zip", req.Name)
	zipFilePath := filepath.Join(tempDir, zipFileName)

	logger.Debug("Creating ZIP file", logger.String("path", zipFilePath))
	if err := createZip(projectDir, zipFilePath); err != nil {
		logger.Error("Failed to create ZIP", logger.Err(err))
		if Metrics != nil {
			Metrics.IncrementCounter("projects_generated_total", map[string]string{"status": "failed"})
		}
		c.JSON(500, types.GenerateResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to create ZIP: %v", err),
		})
		return
	}

	logger.Info("ZIP file created successfully", logger.String("file", zipFileName))

	if Metrics != nil {
		Metrics.IncrementCounter("projects_generated_total", map[string]string{"status": "success"})
	}

	// Send ZIP file
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFileName))
	c.File(zipFilePath)
}

// createZip creates a ZIP archive of the specified directory
func createZip(sourceDir, targetZip string) error {
	zipFile, err := os.Create(targetZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(filepath.Dir(sourceDir), path)
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == sourceDir {
			return nil
		}

		// Create header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(relPath)
		header.Method = zip.Deflate

		if info.IsDir() {
			header.Name += "/"
		}

		// Write header
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		// If not a directory, write file content
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
