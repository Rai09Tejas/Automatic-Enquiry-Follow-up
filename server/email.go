package main

import "fmt"

func logConfirmationEmail(name, email, enquiryType string) {
	fmt.Printf("TO: %s\nSUBJECT: Thanks for your enquiry\nBODY: Hi %s, we've received your %s enquiry and will be in touch shortly.\n",
		email, name, enquiryType)
}
