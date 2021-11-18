package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

var (
	cfgFile string
	rover   string
	camera  string
	rootCmd = &cobra.Command{
		Use: "rover-images [OPTIONS] [COMMANDS]",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Mars Rover Images Query CLI")
		},
	}
	Client = &http.Client{}
)

func init() {
	// initialize the configuration based on the provided, or default, file path. Then, initialize the image cache
	cobra.OnInitialize(initConf, initCache)

	// Add persistant flags to allow for customized content based on user input
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVarP(&rover, "rover", "r", "curiosity", "name of the rover to get images from")
	rootCmd.PersistentFlags().StringVarP(&camera, "camera", "C", "NAVCAM", "name of the camera to get images from")
	viper.BindPFlag("rover-name", rootCmd.PersistentFlags().Lookup("rover"))
	viper.BindPFlag("camera-name", rootCmd.PersistentFlags().Lookup("camera"))
	rootCmd.AddCommand(getCmd)
}

func initConf() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find user's home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search for "config".
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
		viper.SetDefault(CacheFile, home+"/.rover-images.cache")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.Debugln("Using config file:", viper.ConfigFileUsed())
	}
}

func Run() error {
	return rootCmd.Execute()
}
