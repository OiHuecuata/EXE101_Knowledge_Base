package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	_ = gin.Default()
	_ = viper.New()
}
