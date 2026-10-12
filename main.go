package main

import (
	"SDP/subsystemA"
	"SDP/subsystemB"
	"SDP/subsystemC"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const storageDirectory = "./"

func readInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
func main() {
	ra, _ := subsystemA.NewIntakeRepository(storageDirectory)
	rb, _ := subsystemB.NewSeatRepository(storageDirectory)
	rc, _ := subsystemC.NewConfirmationRepository(storageDirectory)
	pa := subsystemA.NewPrerequisiteChecker()
	ha := subsystemA.NewFinancialHoldChecker()
	sb := subsystemB.NewSeatAllocator()
	bc := subsystemC.NewBillingSystem()
	tc := subsystemC.NewTranscriptService()
	studentFacade, err := NewStudentRegistrationFacade(
		storageDirectory, pa, ha, ra, sb, rb, bc, tc, rc,
	)
	if err != nil {
		fmt.Printf("Startup error: %v\n", err)
		return
	}
	journal, _ := NewJournalStore(storageDirectory)
	adminFacade := NewAdminInspectionFacade(ra, rb, rc, journal, studentFacade)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n  UNIVERSITY REGISTRATION")
		fmt.Println("1. [Custom Input] Register student with custom values")
		fmt.Println("2. [Happy Path] Complete A -> B -> C (STUDENT_OK, CS201)")
		fmt.Println("3. [Step 0 Failure] Missing prerequisites (STUDENT_FAIL)")
		fmt.Println("4. [Step B Failure] Course full -> 1-step rollback of A")
		fmt.Println("5. [Step C Failure] Debt limit -> 2-step rollback: B then A")
		fmt.Println("6. [Inspection] Inspect persisted disk state by RegistrationID")
		fmt.Println("7. [Crash Recovery] Scan & recover stalled chains after restart")
		fmt.Println("8. [Reset Store] Wipe all persisted JSON files")
		fmt.Println("0. Exit")
		choice := readInput(reader, "Select an option: ")
		switch choice {
		case "1":
			fmt.Println("\nRegistration Input")
			regID := readInput(reader, "Enter Registration ID (e.g. REG-101): ")
			studentID := readInput(reader, "Enter Student ID (STUDENT_OK, STUDENT_FAIL, BROKE_STUDENT, or custom): ")
			courseID := readInput(reader, "Enter Course ID (CS201, FULL_COURSE, or custom): ")
			credStr := readInput(reader, "Enter Credits (default 3): ")
			credits, err := strconv.Atoi(credStr)
			if err != nil || credits <= 0 {
				credits = 3
			}
			fmt.Println("\nOrchestrating registration chain via Facade...")
			err = studentFacade.RegisterCourse(regID, studentID, courseID, credits)
			if err != nil {
				fmt.Printf("\n[RESULT]: Chain failed and compensated: %v\n", err)
			} else {
				fmt.Println("\n[RESULT]: SUCCESS! Course registered and committed to disk.")
			}
		case "2":
			regID := "DEMO-HAPPY-01"
			_ = studentFacade.RegisterCourse(regID, "STUDENT_OK", "CS201", 3)
		case "3":
			regID := "DEMO-FAIL-A-01"
			_ = studentFacade.RegisterCourse(regID, "STUDENT_FAIL", "CS201", 3)
		case "4":
			regID := "DEMO-FAIL-B-01"
			_ = studentFacade.RegisterCourse(regID, "STUDENT_OK", "FULL_COURSE", 3)
		case "5":
			regID := "DEMO-FAIL-C-01"
			_ = studentFacade.RegisterCourse(regID, "BROKE_STUDENT", "CS201", 3)
		case "6":
			regID := readInput(reader, "Enter Registration ID to inspect: ")
			summary, err := adminFacade.InspectRegistration(regID)
			if err != nil {
				fmt.Println("Error inspecting:", err)
			} else {
				fmt.Printf("\nDISK PERSISTENCE REPORT FOR %s\n", regID)
				fmt.Printf("Journal Progress:        %s\n", summary.JournalState)
				fmt.Printf("Subsystem A (Intake):    %s\n", summary.IntakeStatus)
				fmt.Printf("Subsystem B (Seats):     %s\n", summary.SeatStatus)
				fmt.Printf("Subsystem C (Billing):   %s\n", summary.ConfirmStatus)
			}
		case "7":
			fmt.Println("\nScanning chain_journal.json from disk...")
			rec, _ := adminFacade.RecoverStalledChains()
			if len(rec) == 0 {
				fmt.Println("All chains are consistent. Nothing to recover.")
			} else {
				for _, r := range rec {
					fmt.Printf("Recovered: %s\n", r)
				}
			}
		case "8":
			fmt.Println("\nClearing persisted data files...")
			_ = ra.Clear()
			_ = rb.Clear()
			_ = rc.Clear()
			_ = journal.Clear()
			sb.Reset()
			fmt.Println("[SUCCESS]: All JSON storage files deleted. Seat allocator reset.")
		case "0":
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}
