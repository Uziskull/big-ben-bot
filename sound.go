package main

import (
    "encoding/binary"
    "fmt"
    "io"
    "os"
    "time"
    
    discord "github.com/bwmarrin/discordgo"
)

func loadSound(path string) ([][]byte, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("Error opening dca file: %w", err)
    }

    var opuslen int16
    soundBuffer := make([][]byte, 0)
    for {
        // Read opus frame length from dca file.
        err = binary.Read(file, binary.LittleEndian, &opuslen)

        // If this is the end of the file, just return.
        if err == io.EOF || err == io.ErrUnexpectedEOF {
            err := file.Close()
            if err != nil {
                return nil, err
            }
            return soundBuffer, nil
        }

        if err != nil {
            return nil, fmt.Errorf("Error reading from dca file: %w", err)
        }

        // Read encoded pcm from dca file.
        InBuf := make([]byte, opuslen)
        err = binary.Read(file, binary.LittleEndian, &InBuf)

        // Should not be any end of file errors
        if err != nil {
            return nil, fmt.Errorf("Error reading from dca file: %w", err)
        }

        // Append encoded pcm data to the buffer.
        soundBuffer = append(soundBuffer, InBuf)
    }
}

// playSound plays the current buffer to the provided channel.
func playSound(guildID, channelID string, numRings int) (err error) {

    // Join the provided voice channel.
    vc, err := discordClient.ChannelVoiceJoin(guildID, channelID, false, true)
    if err != nil {
        return err
    }
    defer vc.Disconnect()

    // Sleep for a specified amount of time before playing the sound
    time.Sleep(250 * time.Millisecond)
    
    ringBellInVc(initBellBuffer, vc)

    for _ = range numRings {
        ringBellInVc(mainBellBuffer, vc)
    }

    return nil
}

func ringBellInVc(soundBuffer [][]byte, vc *discord.VoiceConnection) {
    // Start speaking.
    vc.Speaking(true)

    // Send the buffer data.
    for _, buff := range soundBuffer {
        vc.OpusSend <- buff
    }

    // Stop speaking
    vc.Speaking(false)

    // Sleep for a specificed amount of time before ending.
    time.Sleep(250 * time.Millisecond)
}
