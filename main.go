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
    fmt.Print("Enter the start date (YYYY-MM-DD): ")
    reader := bufio.NewReader(os.Stdin)
    startDateInput, _ := reader.ReadString('\n')
    startDateInput = strings.TrimSpace(startDateInput)

    startDate, err := time.Parse("2006-01-02", startDateInput)
    if err != nil {
        fmt.Println("Invalid start date format. Please use YYYY-MM-DD.")
        return
    }

    fmt.Print("Enter the end date (YYYY-MM-DD): ")
    endDateInput, _ := reader.ReadString('\n')
    endDateInput = strings.TrimSpace(endDateInput)

    endDate, err := time.Parse("2006-01-02", endDateInput)
    if err != nil {
        fmt.Println("Invalid end date format. Using tomorrow's date as the end date.")
        endDate = time.Now().AddDate(0, 0, 1)
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
    exerciseSet := make(map[string]struct{})

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

        if !recordDate.Before(startDate) && !recordDate.After(endDate) {
            dateStr := recordDate.Format("2006-01-02")
            exercise := re.ReplaceAllString(record[4], "")

            if _, exists := exercisesByDate[dateStr]; !exists {
                exercisesByDate[dateStr] = make(map[string]int)
            }
            exercisesByDate[dateStr][exercise]++
            exerciseSet[exercise] = struct{}{}
        }
    }

    uniqueExercises := getSortedKeys(exerciseSet)

    outputFile, err := os.Create("output.csv")
    if err != nil {
        fmt.Println("Error creating output file:", err)
        return
    }
    defer outputFile.Close()

    writer := csv.NewWriter(outputFile)
    defer writer.Flush()

    header := append([]string{"Date"}, uniqueExercises...)
    writer.Write(header)

    for date := startDate; !date.After(endDate); date = date.AddDate(0, 0, 1) {
        dateStr := date.Format("2006-01-02")
        row := []string{dateStr}
        for _, exercise := range uniqueExercises {
            count := exercisesByDate[dateStr][exercise]
            row = append(row, fmt.Sprintf("%d", count))
        }
        writer.Write(row)
    }

    fmt.Println("CSV file 'output.csv' has been created.")
}

func getSortedKeys(set map[string]struct{}) []string {
    keys := make([]string, 0, len(set))
    for key := range set {
        keys = append(keys, key)
    }
    sort.Strings(keys)
    return keys
} 