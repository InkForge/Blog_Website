package infrastructures

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gomail.v2"
)

// Mock Dialer
type MockDialer struct {
	mock.Mock
}

func (m *MockDialer) DialAndSend(msgs ...*gomail.Message) error {
	args := m.Called(msgs)
	return args.Error(0)
}

// Suite
type SMTPServiceTestSuite struct {
	suite.Suite
	mockDialer  *MockDialer
	smtpService SMTPService
	emailFrom   string
}

// Setup the suite test
func (s *SMTPServiceTestSuite) SetupTest() {
	s.mockDialer = new(MockDialer)
	s.emailFrom = "noreply@example.com"
	s.smtpService = SMTPService{
		dialer:    s.mockDialer,
		EmailFrom: s.emailFrom,
	}
}

func (s *SMTPServiceTestSuite) TestSendEmail_Success() {
	to := "receiver@example.com"
	subject := "Test Subject"
	body := "<p>Hello!</p>"

	s.mockDialer.On("DialAndSend", mock.Anything).Return(nil).Once()

	err := s.smtpService.SendEmail(to, subject, body)
	s.NoError(err)
	s.mockDialer.AssertCalled(s.T(), "DialAndSend", mock.Anything)
	s.mockDialer.AssertExpectations(s.T())
}

// Run the suite
func TestSMTPServiceSuite(t *testing.T) {
	suite.Run(t, new(SMTPServiceTestSuite))
}
