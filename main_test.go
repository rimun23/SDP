package main

import (
	"SDP/subsystemA"
	"SDP/subsystemB"
	"SDP/subsystemC"
	"os"
	"testing"
)

const testDir = "./test_storage/"

type RegistrationDSL struct {
	t         *testing.T
	facade    *StudentRegistrationFacade
	admin     *AdminInspectionFacade
	studentID string
	courseID  string
	regID     string
	credits   int
	lastErr   error
}

func GivenStudent(t *testing.T, sf *StudentRegistrationFacade, af *AdminInspectionFacade, studentID string) *RegistrationDSL {
	return &RegistrationDSL{
		t:         t,
		facade:    sf,
		admin:     af,
		studentID: studentID,
		credits:   3,
	}
}
func (d *RegistrationDSL) TryingToRegister(regID, courseID string) *RegistrationDSL {
	d.regID = regID
	d.courseID = courseID
	return d
}
func (d *RegistrationDSL) WhenExecutingChain() *RegistrationDSL {
	d.lastErr = d.facade.RegisterCourse(d.regID, d.studentID, d.courseID, d.credits)
	return d
}
func (d *RegistrationDSL) ThenExpectSuccess() *RegistrationDSL {
	if d.lastErr != nil {
		d.t.Fatalf("expected registration success, but got error: %v", d.lastErr)
	}
	if a := d.facade.repoA.Find(d.regID); a == nil || a.Status != "RESERVED" {
		d.t.Errorf("Subsystem A must be RESERVED")
	}
	if b := d.facade.repoB.Find(d.regID); b == nil || b.Status != "ALLOCATED" {
		d.t.Errorf("Subsystem B must be ALLOCATED")
	}
	if c := d.facade.repoC.Find(d.regID); c == nil || c.Status != "CONFIRMED" {
		d.t.Errorf("Subsystem C must be CONFIRMED")
	}
	return d
}
func (d *RegistrationDSL) ThenExpectFailureWithCompensatedState() *RegistrationDSL {
	if d.lastErr == nil {
		d.t.Fatalf("expected chain to fail and compensate, but it succeeded")
	}
	logEntry := d.facade.journal.GetAll()[d.regID]
	if logEntry == nil || logEntry["state"] != "COMPENSATED" {
		d.t.Errorf("expected journal status COMPENSATED, got: %v", logEntry)
	}
	return d
}

func (d *RegistrationDSL) AndAssertSubsystemAIsCancelled() *RegistrationDSL {
	recA := d.facade.repoA.Find(d.regID)
	if recA == nil || recA.Status != "CANCELLED" {
		d.t.Errorf("Subsystem A was not rolled back to CANCELLED")
	}
	return d
}
func (d *RegistrationDSL) AndAssertSubsystemBIsReleased() *RegistrationDSL {
	recB := d.facade.repoB.Find(d.regID)
	if recB == nil || recB.Status != "RELEASED" {
		d.t.Errorf("Subsystem B was not rolled back to RELEASED")
	}
	return d
}
func setupTestEnvironment(t *testing.T) (*StudentRegistrationFacade, *AdminInspectionFacade, func()) {
	_ = os.RemoveAll(testDir)
	_ = os.MkdirAll(testDir, 0755)

	repoA := subsystemA.NewIntakeRepository(testDir)
	repoB := subsystemB.NewSeatRepository(testDir)
	repoC := subsystemC.NewConfirmationRepository(testDir)
	allocB := subsystemB.NewSeatAllocator()
	journal := &JournalStore{path: testDir + "journal.json"}

	sf := &StudentRegistrationFacade{
		checker: &subsystemA.PrerequisiteChecker{},
		repoA:   repoA,
		allocB:  allocB,
		repoB:   repoB,
		billC:   &subsystemC.BillingSystem{},
		repoC:   repoC,
		journal: journal,
	}
	af := &AdminInspectionFacade{repoA: repoA, repoB: repoB, repoC: repoC, journal: journal, student: sf}

	cleanup := func() {
		_ = os.RemoveAll(testDir)
	}
	return sf, af, cleanup
}

func TestHappyPathRegistration(t *testing.T) {
	sf, af, cleanup := setupTestEnvironment(t)
	defer cleanup()

	GivenStudent(t, sf, af, "STUDENT_OK").
		TryingToRegister("TEST-HAPPY", "CS201").
		WhenExecutingChain().
		ThenExpectSuccess()
}

func TestPreCommitFailure_PrerequisiteNotMet(t *testing.T) {
	sf, af, cleanup := setupTestEnvironment(t)
	defer cleanup()

	dsl := GivenStudent(t, sf, af, "STUDENT_FAIL").
		TryingToRegister("TEST-PRECOMMIT", "CS201").
		WhenExecutingChain()

	if dsl.lastErr == nil {
		t.Fatalf("expected prerequisite error, but got nil")
	}
	if recA := sf.repoA.Find("TEST-PRECOMMIT"); recA != nil {
		t.Errorf("Subsystem A should not persist anything on Step 0 failure")
	}
}

func TestDepth1Compensation_StepBFailsRollbackA(t *testing.T) {
	sf, af, cleanup := setupTestEnvironment(t)
	defer cleanup()

	GivenStudent(t, sf, af, "STUDENT_OK").
		TryingToRegister("TEST-B-FAIL", "FULL_COURSE").
		WhenExecutingChain().
		ThenExpectFailureWithCompensatedState().
		AndAssertSubsystemAIsCancelled()
}

func TestDepth2Compensation_StepCFailsRollbackBThenA(t *testing.T) {
	sf, af, cleanup := setupTestEnvironment(t)
	defer cleanup()

	GivenStudent(t, sf, af, "BROKE_STUDENT").
		TryingToRegister("TEST-C-FAIL", "CS201").
		WhenExecutingChain().
		ThenExpectFailureWithCompensatedState().
		AndAssertSubsystemBIsReleased().
		AndAssertSubsystemAIsCancelled()
}

func TestPersistenceReloadRecovery(t *testing.T) {
	sf, af, cleanup := setupTestEnvironment(t)
	defer cleanup()

	regID := "CRASH-TEST-01"
	sf.repoA.Save(&subsystemA.IntakeRecord{RegistrationID: regID, StudentID: "STUDENT_OK", CourseID: "CS201", Status: "RESERVED"})
	sf.repoB.Save(&subsystemB.SeatAllocation{RegistrationID: regID, CourseID: "CS201", SeatNumber: 1, Status: "ALLOCATED"})
	sf.journal.Log(regID, "PENDING_C", "Server crashed before Step C")

	af.Recover()

	if recB := sf.repoB.Find(regID); recB == nil || recB.Status != "RELEASED" {
		t.Errorf("Subsystem B was not recovered to RELEASED")
	}
	if recA := sf.repoA.Find(regID); recA == nil || recA.Status != "CANCELLED" {
		t.Errorf("Subsystem A was not recovered to CANCELLED")
	}
	if logEntry := sf.journal.GetAll()[regID]; logEntry["state"] != "COMPENSATED" {
		t.Errorf("Journal was not updated to COMPENSATED")
	}
}
