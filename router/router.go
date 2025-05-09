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
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/iLogtail/ConfigServer/interface/agent"
	"github.com/iLogtail/ConfigServer/interface/user"
	"github.com/iLogtail/ConfigServer/setting"
)

//go:embed statics
var static embed.FS

func InitRouter() {
	router := gin.Default()

	InitUserRouter(router)
	InitAgentRouter(router)

	router.Any("/api/v1/*proxyPath", reverseProxy(setting.GetSetting().IP+":"+setting.GetSetting().Port))

	router.StaticFS("/", http.FS(static))

	err := router.Run(setting.GetSetting().IP + ":" + setting.GetSetting().Port)
	if err != nil {
		panic(err)
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

// reverseProxy 反向代理
func reverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(target)
		if err != nil {
			c.String(http.StatusInternalServerError, "Bad target")
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/api/v1")
		c.Request.Host = remote.Host
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
