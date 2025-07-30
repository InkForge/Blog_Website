# 🧮 Student Grade Calculator

This is a simple command-line application written in Go that allows users to:

* Enter their name
* Input the number of subjects they studied
* Enter subject names and corresponding grades
* Calculate and display the average grade in a formatted report

---

## 📦 Features

* Validates that the name input only contains alphabetic characters and spaces.
* Ensures that grades are numeric and between 0 and 100.
* Collects subject names and their respective grades.
* Calculates the average grade across all entered subjects.
* Displays a neatly formatted grade report.

---

## 🚀 How to Run

1. **Clone or download the project.**

2. **Run the Go file:**

```bash
go run main.go
```

3. **Follow the interactive prompts:**

```text
***************************************
        Student Grade Calculator

Enter your full name: John Doe
Enter how many subjects you have learned: 3
Enter Subject name: Math
Enter Grade for Math: 90
Enter Subject name: Science
Enter Grade for Science: 80
Enter Subject name: History
Enter Grade for History: 85
```

4. **You’ll get a report like:**

```
******Grade Report for John Doe******
Number of Subjects: 3

#No  Subjects          Grades
1     Math              90
2     Science           80
3     History           85

Average = 85
```

---

## 🔐 Input Validation Rules

* **Name**: must contain only letters and spaces.
* **Subject Names**: must follow the same rule as name.
* **Grades**: must be integers between 0 and 100.
* **Invalid inputs** will re-prompt the user until valid input is given.

---

## 🧱 Code Structure

* `Student` struct: holds user details and grades
* `getStringInput`: reads validated string input using `bufio.Reader`
* `getIntInput`: reads validated integer input (0–100)
* `calculateAverge`: computes average grade
* `display`: prints the grade report

---

## 🛠 Built With

* Go 1.21+
* Standard libraries: `fmt`, `bufio`, `os`, `strings`, `strconv`, `unicode`

---

## 📄 License

Free to use and modify for educational or personal use.

---

