package main

import (
	"fmt"
	"os"

	"github.com/UramnOIL/df-stdcmd/commands"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	conf, err := readConfig(log)
	if err != nil {
		log.Fatalln(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	cmd.Register(cmd.New("kill", "Commit suicide or kill other players", nil, commands.SuicideCommand{}, commands.KillCommand{}))
	cmd.Register(cmd.New("tp", "Teleport entites", []string{"tp"}, commands.TeleportToTargetCommand{}, commands.TeleportToCoordinateCommand{}, commands.TeleportVictimToTargetCommand{}, commands.TeleportVictimToCoordinateCommand{}))
	cmd.Register(cmd.New("gamemode", "Set players' gamemode", nil, commands.SetMyGameModeFromStringCommand{}, commands.SetMyGameModeFromIntCommand{}, commands.SetTargetGameModeFromStringCommand{}, commands.SetTargetGameModeFromIntCommand{}))
	cmd.Register(cmd.New("kick", "kick a player off server", nil, commands.KickCommand{}))

	srv.Listen()
	for srv.Accept(nil) {
	}
}

// readConfig reads the configuration from the config.toml file, or creates the
// file if it does not yet exist.
func readConfig(log server.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return zero, nil
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}
