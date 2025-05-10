// Copyright 2025 iLogtail Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"embed"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/iLogtail/ConfigServer/interface/agent"
	"github.com/iLogtail/ConfigServer/interface/user"
	"github.com/iLogtail/ConfigServer/setting"
)

//go:embed statics/*
var static embed.FS

func InitRouter() {
	router := gin.Default()

	InitUserRouter(router)
	InitAgentRouter(router)

	// rewrite /api/v1/
	router.Any("/api/v1/*proxyPath", func(c *gin.Context) {
		p := c.Param("proxyPath")
		if strings.HasPrefix(p, "/User/") || strings.HasPrefix(p, "/Agent/") {
			c.Request.URL.Path = p
			router.HandleContext(c)
			return
		}
		c.String(http.StatusNotFound, "unknown path: %s", p)
	})

	// 只处理静态文件
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果是 /api/ 路径，返回 404
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API not found"})
			return
		}

		// default return index.html
		if path == "/" || path == "" {
			path = "/index.html"
		}
		embedPath := filepath.Join("statics", strings.TrimPrefix(path, "/"))

		data, err := static.ReadFile(embedPath)
		if err != nil {
			c.String(http.StatusNotFound, "404 not found: %s", path)
			return
		}
		c.Header("Content-Type", detectContentType(path))
		c.Writer.WriteHeader(http.StatusOK)
		_, _ = c.Writer.Write(data)
	})

	err := router.Run(setting.GetSetting().IP + ":" + setting.GetSetting().Port)
	if err != nil {
		panic("Failed to start server: " + err.Error())
	}
}

func InitUserRouter(router *gin.Engine) {
	userGroup := router.Group("/User")
	{
		userGroup.POST("/CreateAgentGroup", user.CreateAgentGroup)
		userGroup.PUT("/UpdateAgentGroup", user.UpdateAgentGroup)
		userGroup.DELETE("/DeleteAgentGroup", user.DeleteAgentGroup)
		userGroup.POST("/GetAgentGroup", user.GetAgentGroup)
		userGroup.POST("/ListAgentGroups", user.ListAgentGroups)

		userGroup.POST("/CreateConfig", user.CreateConfig)
		userGroup.PUT("/UpdateConfig", user.UpdateConfig)
		userGroup.DELETE("/DeleteConfig", user.DeleteConfig)
		userGroup.POST("/GetConfig", user.GetConfig)
		userGroup.POST("/ListConfigs", user.ListConfigs)

		userGroup.PUT("/ApplyConfigToAgentGroup", user.ApplyConfigToAgentGroup)
		userGroup.DELETE("/RemoveConfigFromAgentGroup", user.RemoveConfigFromAgentGroup)
		userGroup.POST("/GetAppliedConfigsForAgentGroup", user.GetAppliedConfigsForAgentGroup)
		userGroup.POST("/GetAppliedAgentGroups", user.GetAppliedAgentGroups)
		userGroup.POST("/ListAgents", user.ListAgents)
	}
}

func InitAgentRouter(router *gin.Engine) {
	agentGroup := router.Group("/Agent")
	{
		agentGroup.POST("/HeartBeat", agent.HeartBeat)

		agentGroup.POST("/FetchPipelineConfig", agent.FetchPipelineConfig)
		agentGroup.POST("/FetchAgentConfig", agent.FetchAgentConfig)
	}
}

// detectContentType detects the content type based on the file extension.
func detectContentType(name string) string {
	switch filepath.Ext(name) {
	case ".html":
		return "text/html; charset=utf-8"
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}
