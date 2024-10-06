package main

import (
    "log/slog"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    discord "github.com/bwmarrin/discordgo"
    "github.com/robfig/cron/v3"
)

var (
    discordClient   *discord.Session
    cronScheduler   *cron.Cron
    
    userVoiceCache  = make(map[string]*UserStatus)
    
    initBellBuffer  = make([][]byte, 0)
    mainBellBuffer  = make([][]byte, 0)
)

type UserStatus struct {
    GuildID   string
    ChannelID string
}

func main() {
    var err error
    start := time.Now()

    timezone := os.Getenv("TZ")
    if timezone == "" {
        timezone = "UTC"
    }
    slog.Info("Today, Big Ben is using", "timezone", timezone)

    initBellBuffer, err = loadSound("initBell.dca")
    if err != nil {
        slog.Error("Error loading init bell sound", "error", err.Error())
        os.Exit(1)
    }

    mainBellBuffer, err = loadSound("mainBell.dca")
    if err != nil {
        slog.Error("Error loading main bell sound", "error", err.Error())
        os.Exit(1)
    }
    
    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        slog.Error("No token provided through env 'BOT_TOKEN'")
        os.Exit(1)
    }
    
    discordClient, err = discord.New("Bot " + token)
    if err != nil {
        slog.Error("Error creating Discord session", "error", err.Error())
        os.Exit(1)
    }
    
    discordClient.AddHandler(readyCallback)
    discordClient.AddHandler(voiceStateCallback)
    discordClient.Identify.Intents = discord.IntentsGuildVoiceStates
    
    err = discordClient.Open()
    if err != nil {
        slog.Error("Error opening Discord session", "error", err.Error())
        os.Exit(1)
    }
    
    slog.Info("Bot is up and running", "startup_time", time.Since(start))
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-sc

    slog.Info("Shutting down...")
    if cronScheduler != nil {
        cronScheduler.Stop()
    }
    discordClient.Close()
}

func readyCallback(s *discord.Session, event *discord.Ready) {
    // need at least one ready guild for crons to start happening
    slog.Info("Bot is ready!")
    if cronScheduler == nil {
        cronScheduler = cron.New()
        cronScheduler.AddFunc("@hourly", playSoundHourly)
        cronScheduler.Start()
    }
}

func voiceStateCallback(s *discord.Session, event *discord.VoiceStateUpdate) {
    slog.Info("Updating cache")

    if event.VoiceState.UserID != "" {
        if _, exists := userVoiceCache[event.VoiceState.UserID]; exists && event.VoiceState.ChannelID == "" {
            delete(userVoiceCache, event.VoiceState.UserID)
            return
        }
        userVoiceCache[event.VoiceState.UserID] = &UserStatus{
            GuildID:   event.VoiceState.GuildID,
            ChannelID: event.VoiceState.ChannelID,
        }
    }
}

func playSoundHourly() {
    currentHour := time.Now().Hour() % 12
    if currentHour == 0 {
        currentHour = 12
    }
    slog.Info("Clock struck the hour - time to ring the bell!", "current_hour", currentHour)
    
    activeVoiceChannels := make(map[string]string)
    for _, userStatus := range userVoiceCache {
        if _, exists := activeVoiceChannels[userStatus.ChannelID]; !exists {
            activeVoiceChannels[userStatus.ChannelID] = userStatus.GuildID
        }
    }

    for chID, gID := range activeVoiceChannels {
        err := playSound(gID, chID, currentHour)
        if err != nil {
            slog.Error("Error playing sound", "error", err.Error())
        }
    }
}
