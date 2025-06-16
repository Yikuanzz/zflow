package server

import (
	"fmt"
	"net/http"

	"zflow/app/zflow/data"
	"zflow/app/zflow/model"

	"github.com/gin-gonic/gin"
)

func NewServer() *http.Server {
	router := gin.Default()

	router.GET("/node_types", func(c *gin.Context) {
		// TODO: 从注册中心获取服务地址

		// TODO: 根据服务地址获取所有的节点类型

		// TODO: 返回所有的节点类型与本地节点类型合并
		c.JSON(http.StatusOK, data.GetDefaultNodeTypes())
	})

	router.GET("/connection_types", func(c *gin.Context) {
		// TODO: 从注册中心获取服务地址

		// TODO: 根据服务地址获取所有的连接类型

		// TODO: 返回所有的连接类型与本地连接类型合并
		c.JSON(http.StatusOK, data.GetDefaultConnectionTypes())
	})

	router.POST("/workflows", func(c *gin.Context) {
		var req struct {
			UID      string            `json:"uid"`
			Workflow model.RawWorkflow `json:"workflow"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.UID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "uid is required"})
			return
		}
		if req.Workflow.Nodes == nil || req.Workflow.Connections == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "workflow is required"})
			return
		}

		// 1、创建工作流
		wf, err := model.NewWorkflow(req.UID, req.Workflow)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2、注入节点和连接类型的操作
		if err := data.InjectOperations(wf); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 3、创建执行上下文
		ctx := &model.ExecutionContext{
			Workflow: wf,
			Logger: func(msg string) {
				fmt.Println(msg)
			},
			Vars: make(map[string]interface{}),
		}

		// 4、执行工作流
		if err := wf.ExecuteWorkflow(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 5、收集工作流执行结果
		result := wf.CollectWorkflowResults()

		c.JSON(http.StatusOK, result)
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}
