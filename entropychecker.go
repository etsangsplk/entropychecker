package entropychecker

import (
	"errors"
	"ioutil"
	"strconv"
	"time"
)

// Set this to what you consider to be a 'safe' minimum entropy amount (in bits)
var EntropyLimit = 128

// Waiting for entropy will time out after this amount of time. Setting to zero will never time out.
var Timeout = time.Second * 10

var ErrTimeout = errors.New("Timed out waiting for sufficient entropy")

// Get the entropy estimate. Returns the estimated entropy in bits
func GetEntropy() (int, error) {
	text, err := ioutil.ReadFile("/proc/sys/kernel/random/entropy_avail")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(text)
}

// Block until sufficient entropy is available
func WaitForEntropy() error {
	// set up the timeout
	timeout := make(chan bool, 1)
	if Timeout != 0 {
		go func() {
			time.Sleep(Timeout)
			timeout <- true
		}()
	}

	for {
		entropy, err := GetEntropy()

		switch {
		case err != nil:
			close(timeout)
			return err
		case entropy > EntropyLimit:
			close(timeout)
			return nil
		default:
			select {
			case <-timeout:
				return ErrTimeout
			default:
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
}
