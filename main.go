package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	var timeArg string
	var message string
	var iconPath string
	var expireTime int
	var seconds int
	var delay time.Duration

	// Root command
	var rootCmd = &cobra.Command{
		Use:   "notify",
		Short: "Schedule a notification",
    Example: "notify -s 5 -m 'Hello, World!'\nnotify -t '21:00' -m 'Good evening'",
    Version: "0.0.1",
		Run: func(cmd *cobra.Command, args []string) {
			// Validate the input
			if timeArg == "" && seconds == 0 || message == "" {
				fmt.Println("Error: Must set message and time OR seconds")
				cmd.Usage()
				os.Exit(1)
			}

			if timeArg != "" && seconds > 0 {
				fmt.Println("Error: Cannot set both time and seconds")
				cmd.Usage()
				os.Exit(1)
			}

			if timeArg != "" {
				// Parse the specified time
				now := time.Now()
				notifyTime, err := time.Parse("15:04", timeArg)
				if err != nil {
					fmt.Println("Error: Invalid time format. Use HH:MM")
					os.Exit(1)
				}

				// Adjust the date of notifyTime to today
				notifyTime = time.Date(now.Year(), now.Month(), now.Day(), notifyTime.Hour(), notifyTime.Minute(), 0, 0, now.Location())

				// Calculate the delay
				delay = notifyTime.Sub(now)

				// If the time is in the past, schedule it for tomorrow
				if delay < 0 {
					notifyTime = notifyTime.Add(24 * time.Hour)
					delay = notifyTime.Sub(now)
				}

				fmt.Printf("Notification scheduled for %s\n", notifyTime.Format("15:04"))
			}

			if seconds > 0 {
        fmt.Printf("Seconds %d\n", seconds)
				delay = time.Duration(seconds) * time.Second
				fmt.Printf("Notification scheduled in %.0f seconds\n", delay.Seconds())
			}

			var wg sync.WaitGroup

			wg.Add(1)
      
			// Run the notification scheduling in a separate goroutine
			go func() {
				defer wg.Done()

				if delay == 0 {
					fmt.Println("Error: Delay is not set. Have you set a time or seconds flag?")
					wg.Done()
					return
				}

				time.Sleep(delay)
				cmdArgs := []string{}
				if iconPath != "" {
					cmdArgs = append(cmdArgs, "-i", iconPath)
				}

				if expireTime > 0 {
					cmdArgs = append(cmdArgs, "-t", fmt.Sprintf("%d", expireTime*1000))
				}

				cmdArgs = append(cmdArgs, message)

				cmd := exec.Command("notify-send", cmdArgs...)
				// cmd.Wait()

				err := cmd.Run()
				if err != nil {
					fmt.Println("Error sending notification:", err)
				}
			}()

			wg.Wait()
		},
	}

	// Define flags
	rootCmd.Flags().StringVarP(&timeArg, "time", "t", "", "The time to send the notification (HH:MM)")
	rootCmd.Flags().StringVarP(&message, "message", "m", "", "The message to display in the notification")
	rootCmd.Flags().StringVarP(&iconPath, "icon", "i", "", "The path to an icon to display with the notification")
	rootCmd.Flags().IntVarP(&expireTime, "seconds", "s", 0, "The number of seconds to display the notification. Cannot be set with time")
	rootCmd.Flags().IntVarP(&seconds, "expire-time", "e", 3, "The number of seconds to expire the notification. Defaults to 3")

	// Execute the command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
