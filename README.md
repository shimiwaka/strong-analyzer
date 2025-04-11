# Strong Analyzer

## Overview

Strong Analyzer is a Go application designed to analyze workout data exported from the Strong Workout Tracker app. It processes CSV files to generate a summary of exercises performed over a specified date range.

## Installation

1. Ensure you have Go installed on your system.
2. Clone this repository to your local machine.
3. Navigate to the project directory.

## Usage

1. Place your `strong.csv` file exported from Strongin the project directory.
2. Run the application using the following command:
   ```bash
   go run main.go
   ```
3. Enter the start and end dates in the format `YYYY-MM-DD` when prompted.
4. The application will generate an `output.csv` file with the analysis results.

## CSV Output

The output CSV file will contain:
- A header row with the date and each exercise type.
- Rows for each date in the specified range, showing the count of each exercise performed.
