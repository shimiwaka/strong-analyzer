package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "time"
    "bufio"
    "strings"
    "sort"
    "regexp"
)

func main() {
    fmt.Print("Enter the date (YYYY-MM-DD): ")
    reader := bufio.NewReader(os.Stdin)
    dateInput, _ := reader.ReadString('\n')
    dateInput = strings.TrimSpace(dateInput)

    inputDate, err := time.Parse("2006-01-02", dateInput)
    if err != nil {
        fmt.Println("Invalid date format. Please use YYYY-MM-DD.")
        return
    }

    file, err := os.Open("strong.csv")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    csvReader := csv.NewReader(file)
    csvReader.LazyQuotes = true
    csvReader.Comma = ';'

    records, err := csvReader.ReadAll()
    if err != nil {
        fmt.Println("Error reading CSV:", err)
        return
    }

    if len(records) > 0 && records[0][0] == "Workout #" {
        records = records[1:]
    }

    exercisesByDate := make(map[string]map[string]int)
    var dates []string

    re := regexp.MustCompile(`\s\(.*?\)`)

    for _, record := range records {
        if record[5] == "Rest Timer" {
            continue
        }

        recordDate, err := time.Parse("2006-01-02 15:04:05", record[1])
        if err != nil {
            fmt.Println("Error parsing date in record:", err)
            continue
        }

        if !recordDate.Before(inputDate) {
            dateStr := recordDate.Format("2006-01-02")
            exercise := re.ReplaceAllString(record[4], "")

            if _, exists := exercisesByDate[dateStr]; !exists {
                exercisesByDate[dateStr] = make(map[string]int)
                dates = append(dates, dateStr)
            }
            exercisesByDate[dateStr][exercise]++

            fmt.Printf("Processing record: %v\n", record)
            fmt.Printf("Parsed date: %s\n", recordDate.Format("2006-01-02"))
            fmt.Printf("Exercises for date %s: %v\n", dateStr, exercisesByDate[dateStr])
        }
    }

    dateSet := make(map[string]struct{})
    uniqueDates := []string{}
    for _, date := range dates {
        if _, exists := dateSet[date]; !exists {
            dateSet[date] = struct{}{}
            uniqueDates = append(uniqueDates, date)
        }
    }

    sort.Strings(uniqueDates)

    for _, date := range uniqueDates {
        fmt.Printf("Date: %s, Exercises: ", date)
        for exercise, count := range exercisesByDate[date] {
            fmt.Printf("%s (%d) ", exercise, count)
        }
        fmt.Println()
    }
} 