package main

import (
	"SDP/subsystemA"
	"SDP/subsystemB"
	"SDP/subsystemC"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	dir := "./"
	repoA := subsystemA.NewIntakeRepository(dir)
	repoB := subsystemB.NewSeatRepository(dir)
	repoC := subsystemC.NewConfirmationRepository(dir)
	allocB := subsystemB.NewSeatAllocator()
	journal := &JournalStore{path: dir + "journal.json"}

	studentFac := &StudentRegistrationFacade{
		checker: &subsystemA.PrerequisiteChecker{},
		repoA:   repoA,
		allocB:  allocB,
		repoB:   repoB,
		billC:   &subsystemC.BillingSystem{},
		repoC:   repoC,
		journal: journal,
	}
	adminFac := &AdminInspectionFacade{repoA: repoA, repoB: repoB, repoC: repoC, journal: journal, student: studentFac}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nUNIVERSITY REGISTRATION")
		fmt.Println("1. [Happy Path] (STUDENT_OK, CS201) -> A, B, C Commit")
		fmt.Println("2. [Fail at A]  (STUDENT_FAIL) -> Rejects before commit")
		fmt.Println("3. [Fail at B]  (FULL_COURSE)  -> 1-step rollback of A")
		fmt.Println("4. [Fail at C]  (BROKE_STUDENT)-> 2-step rollback: B then A")
		fmt.Println("5. [Inspect] View persisted state by ID")
		fmt.Println("6. [Recover] Fix broken chains after reboot")
		fmt.Println("7. [Reset] Wipe JSON stores & memory")
		fmt.Println("0. Exit")
		fmt.Print("Choice: ")

		text, _ := reader.ReadString('\n')
		switch strings.TrimSpace(text) {
		case "1":
			_ = studentFac.RegisterCourse("REG-01", "STUDENT_OK", "CS201", 3)
		case "2":
			_ = studentFac.RegisterCourse("REG-02", "STUDENT_FAIL", "CS201", 3)
		case "3":
			_ = studentFac.RegisterCourse("REG-03", "STUDENT_OK", "FULL_COURSE", 3)
		case "4":
			_ = studentFac.RegisterCourse("REG-04", "BROKE_STUDENT", "CS201", 3)
		case "5":
			fmt.Print("Enter ID to inspect (e.g. REG-01, REG-03, REG-04): ")
			id, _ := reader.ReadString('\n')
			adminFac.Inspect(strings.TrimSpace(id))
		case "6":
			adminFac.Recover()
		case "7":
			repoA.Clear()
			repoB.Clear()
			repoC.Clear()
			journal.Clear()
			allocB.Reset()
			fmt.Println("All JSON files deleted and memory reset.")
		case "0":
			return
		}
	}
}
