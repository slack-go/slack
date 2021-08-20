package main

import (
	"encoding/json"
	"fmt"

	"github.com/slack-go/slack"
)

// The functions below mock the different templates slack has as examples on their website.
//
// Refer to README.md for more information on the examples and how to use them.

func main() {

	fmt.Println("--- Begin Example One ---")
	exampleOne()
	fmt.Println("--- End Example One ---")

	fmt.Println("--- Begin Example Two ---")
	exampleTwo()
	fmt.Println("--- End Example Two ---")

	fmt.Println("--- Begin Example Three ---")
	exampleThree()
	fmt.Println("--- End Example Three ---")

	fmt.Println("--- Begin Example Four ---")
	exampleFour()
	fmt.Println("--- End Example Four ---")

	fmt.Println("--- Begin Example Five ---")
	exampleFive()
	fmt.Println("--- End Example Five ---")

	fmt.Println("--- Begin Example Six ---")
	exampleSix()
	fmt.Println("--- End Example Six ---")

	fmt.Println("--- Begin Example Unmarshalling ---")
	unmarshalExample()
	fmt.Println("--- End Example Unmarshalling ---")
}

// approvalRequest mocks the simple "Approval" template located on block kit builder website
func exampleOne() {

	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "You have a new request:\n*<fakeLink.toEmployeeProfile.com|Fred Enriquez - New device request>*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Fields
	typeField := slack.NewTextBlockObject("mrkdwn", "*Type:*\nComputer (laptop)", false, false)
	whenField := slack.NewTextBlockObject("mrkdwn", "*When:*\nSubmitted Aut 10", false, false)
	lastUpdateField := slack.NewTextBlockObject("mrkdwn", "*Last Update:*\nMar 10, 2015 (3 years, 5 months)", false, false)
	reasonField := slack.NewTextBlockObject("mrkdwn", "*Reason:*\nAll vowel keys aren't working.", false, false)
	specsField := slack.NewTextBlockObject("mrkdwn", "*Specs:*\n\"Cheetah Pro 15\" - Fast, really fast\"", false, false)

	fieldSlice := make([]*slack.TextBlockObject, 0)
	fieldSlice = append(fieldSlice, typeField)
	fieldSlice = append(fieldSlice, whenField)
	fieldSlice = append(fieldSlice, lastUpdateField)
	fieldSlice = append(fieldSlice, reasonField)
	fieldSlice = append(fieldSlice, specsField)

	fieldsSection := slack.NewSectionBlock(nil, fieldSlice, nil)

	// Approve and Deny Buttons
	approveBtnTxt := slack.NewTextBlockObject("plain_text", "Approve", false, false)
	approveBtn := slack.NewButtonBlockElement("", "click_me_123", approveBtnTxt)

	denyBtnTxt := slack.NewTextBlockObject("plain_text", "Deny", false, false)
	denyBtn := slack.NewButtonBlockElement("", "click_me_123", denyBtnTxt)

	actionBlock := slack.NewActionBlock("", approveBtn, denyBtn)

	// Build Message with blocks created above

	msg := slack.NewBlockMessage(
		headerSection,
		fieldsSection,
		actionBlock,
	)

	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))

}

// exampleTwo mocks the more complex "Approval" template located on block kit builder website
// which includes an accessory image next to the approval request
func exampleTwo() {

	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "You have a new request:\n*<google.com|Fred Enriquez - Time Off request>*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	approvalText := slack.NewTextBlockObject("mrkdwn", "*Type:*\nPaid time off\n*When:*\nAug 10-Aug 13\n*Hours:* 16.0 (2 days)\n*Remaining balance:* 32.0 hours (4 days)\n*Comments:* \"Family in town, going camping!\"", false, false)
	approvalImage := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/approvalsNewDevice.png", "computer thumbnail")

	fieldsSection := slack.NewSectionBlock(approvalText, nil, slack.NewAccessory(approvalImage))

	// Approve and Deny Buttons
	approveBtnTxt := slack.NewTextBlockObject("plain_text", "Approve", false, false)
	approveBtn := slack.NewButtonBlockElement("", "click_me_123", approveBtnTxt)

	denyBtnTxt := slack.NewTextBlockObject("plain_text", "Deny", false, false)
	denyBtn := slack.NewButtonBlockElement("", "click_me_123", denyBtnTxt)

	actionBlock := slack.NewActionBlock("", approveBtn, denyBtn)

	// Build Message with blocks created above

	msg := slack.NewBlockMessage(
		headerSection,
		fieldsSection,
		actionBlock,
	)

	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))

}

// exampleThree generates the notification example from the block kit builder website
func exampleThree() {

	// Shared Assets for example
	chooseBtnText := slack.NewTextBlockObject("plain_text", "Choose", true, false)
	chooseBtnEle := slack.NewButtonBlockElement("", "click_me_123", chooseBtnText)
	divSection := slack.NewDividerBlock()

	// Header Section
	headerText := slack.NewTextBlockObject("plain_text", "Looks like you have a scheduling conflict with this event:", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Schedule Info Section
	scheduleText := slack.NewTextBlockObject("mrkdwn", "*<fakeLink.toUserProfiles.com|Iris / Zelda 1-1>*\nTuesday, January 21 4:00-4:30pm\nBuilding 2 - Havarti Cheese (3)\n2 guests", false, false)
	scheduleAccessory := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/notifications.png", "calendar thumbnail")
	schedeuleSection := slack.NewSectionBlock(scheduleText, nil, slack.NewAccessory(scheduleAccessory))

	// Conflict Section
	conflictImage := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/notificationsWarningIcon.png", "notifications warning icon")
	conflictText := slack.NewTextBlockObject("mrkdwn", "*Conflicts with Team Huddle: 4:15-4:30pm*", false, false)

	conflictSection := slack.NewContextBlock(
		"",
		[]slack.MixedElement{conflictImage, conflictText}...,
	)

	// Proposese Text
	proposeText := slack.NewTextBlockObject("mrkdwn", "*Propose a new time:*", false, false)
	proposeSection := slack.NewSectionBlock(proposeText, nil, nil)

	// Option 1
	optionOneText := slack.NewTextBlockObject("mrkdwn", "*Today - 4:30-5pm*\nEveryone is available: @iris, @zelda", false, false)
	optionOneSection := slack.NewSectionBlock(optionOneText, nil, slack.NewAccessory(chooseBtnEle))

	// Option 2
	optionTwoText := slack.NewTextBlockObject("mrkdwn", "*Tomorrow - 4-4:30pm*\nEveryone is available: @iris, @zelda", false, false)
	optionTwoSection := slack.NewSectionBlock(optionTwoText, nil, slack.NewAccessory(chooseBtnEle))

	// Option 3
	optionThreeText := slack.NewTextBlockObject("mrkdwn", "*Tomorrow - 6-6:30pm*\nSome people aren't available: @iris, ~@zelda~", false, false)
	optionThreeSection := slack.NewSectionBlock(optionThreeText, nil, slack.NewAccessory(chooseBtnEle))

	// Show More Times Link
	showMoreText := slack.NewTextBlockObject("mrkdwn", "*<fakelink.ToMoreTimes.com|Show more times>*", false, false)
	showMoreSection := slack.NewSectionBlock(showMoreText, nil, nil)

	// Build Message with blocks created above
	msg := slack.NewBlockMessage(
		headerSection,
		divSection,
		schedeuleSection,
		conflictSection,
		divSection,
		proposeSection,
		optionOneSection,
		optionTwoSection,
		optionThreeSection,
		showMoreSection,
	)

	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))

}

// exampleFour profiles a poll example block
func exampleFour() {

	// Shared Assets for example
	divSection := slack.NewDividerBlock()
	voteBtnText := slack.NewTextBlockObject("plain_text", "Vote", true, false)
	voteBtnEle := slack.NewButtonBlockElement("", "click_me_123", voteBtnText)
	profileOne := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/profile_1.png", "Michael Scott")
	profileTwo := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/profile_2.png", "Dwight Schrute")
	profileThree := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/profile_3.png", "Pam Beasely")
	profileFour := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/profile_4.png", "Angela")

	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "*Where should we order lunch from?* Poll by <fakeLink.toUser.com|Mark>", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Option One Info
	optOneText := slack.NewTextBlockObject("mrkdwn", ":sushi: *Ace Wasabi Rock-n-Roll Sushi Bar*\nThe best landlocked sushi restaurant.", false, false)
	optOneSection := slack.NewSectionBlock(optOneText, nil, slack.NewAccessory(voteBtnEle))

	// Option One Votes
	optOneVoteText := slack.NewTextBlockObject("plain_text", "3 votes", true, false)
	optOneContext := slack.NewContextBlock("", []slack.MixedElement{profileOne, profileTwo, profileThree, optOneVoteText}...)

	// Option Two Info
	optTwoText := slack.NewTextBlockObject("mrkdwn", ":hamburger: *Super Hungryman Hamburgers*\nOnly for the hungriest of the hungry.", false, false)
	optTwoSection := slack.NewSectionBlock(optTwoText, nil, slack.NewAccessory(voteBtnEle))

	// Option Two Votes
	optTwoVoteText := slack.NewTextBlockObject("plain_text", "2 votes", true, false)
	optTwoContext := slack.NewContextBlock("", []slack.MixedElement{profileFour, profileTwo, optTwoVoteText}...)

	// Option Three Info
	optThreeText := slack.NewTextBlockObject("mrkdwn", ":ramen: *Kagawa-Ya Udon Noodle Shop*\nDo you like to shop for noodles? We have noodles.", false, false)
	optThreeSection := slack.NewSectionBlock(optThreeText, nil, slack.NewAccessory(voteBtnEle))

	// Option Three Votes
	optThreeVoteText := slack.NewTextBlockObject("plain_text", "No votes", true, false)
	optThreeContext := slack.NewContextBlock("", []slack.MixedElement{optThreeVoteText}...)

	// Suggestions Action
	btnTxt := slack.NewTextBlockObject("plain_text", "Add a suggestion", false, false)
	nextBtn := slack.NewButtonBlockElement("", "click_me_123", btnTxt)
	actionBlock := slack.NewActionBlock("", nextBtn)

	// Build Message with blocks created above
	msg := slack.NewBlockMessage(
		headerSection,
		divSection,
		optOneSection,
		optOneContext,
		optTwoSection,
		optTwoContext,
		optThreeSection,
		optThreeContext,
		divSection,
		actionBlock,
	)

	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))

}

func exampleFive() {

	// Build Header Section Block, includes text and overflow menu

	headerText := slack.NewTextBlockObject("mrkdwn", "We found *205 Hotels* in New Orleans, LA from *12/14 to 12/17*", false, false)

	// Build Text Objects associated with each option
	overflowOptionTextOne := slack.NewTextBlockObject("plain_text", "Option One", false, false)
	overflowOptionTextTwo := slack.NewTextBlockObject("plain_text", "Option Two", false, false)
	overflowOptionTextThree := slack.NewTextBlockObject("plain_text", "Option Three", false, false)

	// Build each option, providing a value for the option
	overflowOptionOne := slack.NewOptionBlockObject("value-0", overflowOptionTextOne, nil)
	overflowOptionTwo := slack.NewOptionBlockObject("value-1", overflowOptionTextTwo, nil)
	overflowOptionThree := slack.NewOptionBlockObject("value-2", overflowOptionTextThree, nil)

	// Build overflow section
	overflow := slack.NewOverflowBlockElement("", overflowOptionOne, overflowOptionTwo, overflowOptionThree)

	// Create the header section
	headerSection := slack.NewSectionBlock(headerText, nil, slack.NewAccessory(overflow))

	// Shared Divider
	divSection := slack.NewDividerBlock()

	// Shared Objects
	locationPinImage := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png", "Location Pin Icon")

	// First Hotel Listing
	hotelOneInfo := slack.NewTextBlockObject("mrkdwn", "*<fakeLink.toHotelPage.com|Windsor Court Hotel>*\n★★★★★\n$340 per night\nRated: 9.4 - Excellent", false, false)
	hotelOneImage := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/tripAgent_1.png", "Windsor Court Hotel thumbnail")
	hotelOneLoc := slack.NewTextBlockObject("plain_text", "Location: Central Business District", true, false)

	hotelOneSection := slack.NewSectionBlock(hotelOneInfo, nil, slack.NewAccessory(hotelOneImage))
	hotelOneContext := slack.NewContextBlock("", []slack.MixedElement{locationPinImage, hotelOneLoc}...)

	// Second Hotel Listing
	hotelTwoInfo := slack.NewTextBlockObject("mrkdwn", "*<fakeLink.toHotelPage.com|The Ritz-Carlton New Orleans>*\n★★★★★\n$340 per night\nRated: 9.1 - Excellent", false, false)
	hotelTwoImage := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/tripAgent_2.png", "Ritz-Carlton New Orleans thumbnail")
	hotelTwoLoc := slack.NewTextBlockObject("plain_text", "Location: French Quarter", true, false)

	hotelTwoSection := slack.NewSectionBlock(hotelTwoInfo, nil, slack.NewAccessory(hotelTwoImage))
	hotelTwoContext := slack.NewContextBlock("", []slack.MixedElement{locationPinImage, hotelTwoLoc}...)

	// Third Hotel Listing
	hotelThreeInfo := slack.NewTextBlockObject("mrkdwn", "*<fakeLink.toHotelPage.com|Omni Royal Orleans Hotel>*\n★★★★★\n$419 per night\nRated: 8.8 - Excellent", false, false)
	hotelThreeImage := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/tripAgent_3.png", "https://api.slack.com/img/blocks/bkb_template_images/tripAgent_3.png")
	hotelThreeLoc := slack.NewTextBlockObject("plain_text", "Location: French Quarter", true, false)

	hotelThreeSection := slack.NewSectionBlock(hotelThreeInfo, nil, slack.NewAccessory(hotelThreeImage))
	hotelThreeContext := slack.NewContextBlock("", []slack.MixedElement{locationPinImage, hotelThreeLoc}...)

	// Action button
	btnTxt := slack.NewTextBlockObject("plain_text", "Next 2 Results", false, false)
	nextBtn := slack.NewButtonBlockElement("", "click_me_123", btnTxt)
	actionBlock := slack.NewActionBlock("", nextBtn)

	// Build Message with blocks created above
	msg := slack.NewBlockMessage(
		headerSection,
		divSection,
		hotelOneSection,
		hotelOneContext,
		divSection,
		hotelTwoSection,
		hotelTwoContext,
		divSection,
		hotelThreeSection,
		hotelThreeContext,
		divSection,
		actionBlock,
	)

	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))

}

func exampleSix() {

	// Shared Assets for example
	divSection := slack.NewDividerBlock()

	// Shared Available Options
	manageTxt := slack.NewTextBlockObject("plain_text", "Manage", true, false)
	editTxt := slack.NewTextBlockObject("plain_text", "Edit it", false, false)
	readTxt := slack.NewTextBlockObject("plain_text", "Read it", false, false)
	saveTxt := slack.NewTextBlockObject("plain_text", "Save it", false, false)

	editOpt := slack.NewOptionBlockObject("value-0", editTxt, nil)
	readOpt := slack.NewOptionBlockObject("value-1", readTxt, nil)
	saveOpt := slack.NewOptionBlockObject("value-2", saveTxt, nil)

	availableOption := slack.NewOptionsSelectBlockElement("static_select", manageTxt, "", editOpt, readOpt, saveOpt)

	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", ":mag: Search results for *Cata*", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Result One
	resultOneTxt := slack.NewTextBlockObject("mrkdwn", "*<fakeLink.toYourApp.com|Use Case Catalogue>*\nUse Case Catalogue for the following departments/roles...", false, false)
	resultOneSection := slack.NewSectionBlock(resultOneTxt, nil, slack.NewAccessory(availableOption))

	// Result Two
	resultTwoTxt := slack.NewTextBlockObject("mrkdwn", "*<fakeLink.toYourApp.com|Customer Support - Workflow Diagram Catalogue>*\nThis resource was put together by members of...", false, false)
	resultTwoSection := slack.NewSectionBlock(resultTwoTxt, nil, slack.NewAccessory(availableOption))

	// Result Three
	resultThreeTxt := slack.NewTextBlockObject("mrkdwn", "*<fakeLink.toYourApp.com|Self-Serve Learning Options Catalogue>*\nSee the learning and development options we...", false, false)
	resultThreeSection := slack.NewSectionBlock(resultThreeTxt, nil, slack.NewAccessory(availableOption))

	// Result Four
	resultFourTxt := slack.NewTextBlockObject("mrkdwn", "*<fakeLink.toYourApp.com|Use Case Catalogue - CF Presentation - [June 12, 2018]>*\nThis is presentation will continue to be updated as...", false, false)
	resultFourSection := slack.NewSectionBlock(resultFourTxt, nil, slack.NewAccessory(availableOption))

	// Result Five
	resultFiveTxt := slack.NewTextBlockObject("mrkdwn", "*<fakeLink.toYourApp.com|Comprehensive Benefits Catalogue - 2019>*\nInformation about all the benfits we offer is...", false, false)
	resultFiveSection := slack.NewSectionBlock(resultFiveTxt, nil, slack.NewAccessory(availableOption))

	// Next Results Button
	// Suggestions Action
	btnTxt := slack.NewTextBlockObject("plain_text", "Next 5 Results", false, false)
	nextBtn := slack.NewButtonBlockElement("", "click_me_123", btnTxt)
	actionBlock := slack.NewActionBlock("", nextBtn)

	// Build Message with blocks created above
	msg := slack.NewBlockMessage(
		headerSection,
		divSection,
		resultOneSection,
		resultTwoSection,
		resultThreeSection,
		resultFourSection,
		resultFiveSection,
		divSection,
		actionBlock,
	)

	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))

}

func unmarshalExample() {
	var msgBlocks []slack.Block

	// Append ActionBlock for marshalling
	btnTxt := slack.NewTextBlockObject("plain_text", "Add a suggestion", false, false)
	nextBtn := slack.NewButtonBlockElement("", "click_me_123", btnTxt)
	approveBtnTxt := slack.NewTextBlockObject("plain_text", "Approve", false, false)
	approveBtn := slack.NewButtonBlockElement("", "click_me_123", approveBtnTxt)
	msgBlocks = append(msgBlocks, slack.NewActionBlock("", nextBtn, approveBtn))

	// Append ContextBlock for marshalling
	profileOne := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/profile_1.png", "Michael Scott")
	profileTwo := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/profile_2.png", "Dwight Schrute")
	textBlockObj := slack.NewTextBlockObject("mrkdwn", "*<fakeLink.toHotelPage.com|Omni Royal Orleans Hotel>*\n★★★★★\n$419 per night\nRated: 8.8 - Excellent", false, false)
	msgBlocks = append(msgBlocks, slack.NewContextBlock("", []slack.MixedElement{profileOne, profileTwo, textBlockObj}...))

	// Append ImageBlock for marshalling
	msgBlocks = append(msgBlocks, slack.NewImageBlock("https://api.slack.com/img/blocks/bkb_template_images/profile_2.png", "some profile", "image-block", textBlockObj))

	// Append DividerBlock for marshalling
	msgBlocks = append(msgBlocks, slack.NewDividerBlock())

	// Append SectionBlock for marshalling
	approvalText := slack.NewTextBlockObject("mrkdwn", "*Type:*\nPaid time off\n*When:*\nAug 10-Aug 13\n*Hours:* 16.0 (2 days)\n*Remaining balance:* 32.0 hours (4 days)\n*Comments:* \"Family in town, going camping!\"", false, false)
	approvalImage := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/approvalsNewDevice.png", "computer thumbnail")
	msgBlocks = append(msgBlocks, slack.NewSectionBlock(approvalText, nil, slack.NewAccessory(approvalImage)), nil)

	// Build Message with blocks created above
	msg := slack.NewBlockMessage(msgBlocks...)

	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))

	// Unmarshal message
	m := slack.Message{}
	if err := json.Unmarshal(b, &m); err != nil {
		fmt.Println(err)
		return
	}

	var respBlocks []slack.Block
	for _, block := range m.Blocks.BlockSet {
		// Need to implement a type switch to determine Block type since the
		// response from Slack could include any/all types under "blocks" key
		switch block.BlockType() {
		case slack.MBTContext:
			var respMixedElements []slack.MixedElement
			contextElements := block.(*slack.ContextBlock).ContextElements.Elements
			// Need to implement a type switch for ContextElements for same reason as Blocks
			for _, elem := range contextElements {
				switch elem.MixedElementType() {
				case slack.MixedElementImage:
					// Assert the block's type to manipulate/extract values
					imageBlockElem := elem.(*slack.ImageBlockElement)
					imageBlockElem.ImageURL = "https://api.slack.com/img/blocks/bkb_template_images/profile_1.png"
					imageBlockElem.AltText = "MichaelScott"
					respMixedElements = append(respMixedElements, imageBlockElem)
				case slack.MixedElementText:
					textBlockElem := elem.(*slack.TextBlockObject)
					textBlockElem.Text = "go go go go go"
					respMixedElements = append(respMixedElements, textBlockElem)
				}
			}
			respBlocks = append(respBlocks, slack.NewContextBlock("new block", respMixedElements...))
		case slack.MBTAction:
			actionBlock := block.(*slack.ActionBlock)
			// Need to implement a type switch for BlockElements for same reason as Blocks
			for _, elem := range actionBlock.Elements.ElementSet {
				switch elem.ElementType() {
				case slack.METImage:
					imageElem := elem.(*slack.ImageBlockElement)
					fmt.Printf("do something with image block element: %v\n", imageElem)
				case slack.METButton:
					buttonElem := elem.(*slack.ButtonBlockElement)
					fmt.Printf("do something with button block element: %v\n", buttonElem)
				case slack.METOverflow:
					overflowElem := elem.(*slack.OverflowBlockElement)
					fmt.Printf("do something with overflow block element: %v\n", overflowElem)
				case slack.METDatepicker:
					datepickerElem := elem.(*slack.DatePickerBlockElement)
					fmt.Printf("do something with datepicker block element: %v\n", datepickerElem)
				case slack.METTimepicker:
					timepickerElem := elem.(*slack.TimePickerBlockElement)
					fmt.Printf("do something with timepicker block element: %v\n", timepickerElem)
				}
			}
			respBlocks = append(respBlocks, block)
		case slack.MBTImage:
			// Simply re-append the block if you want to include it in the response
			respBlocks = append(respBlocks, block)
		case slack.MBTSection:
			respBlocks = append(respBlocks, block)
		case slack.MBTDivider:
			respBlocks = append(respBlocks, block)
		}
	}

	// Build new Message with Blocks obtained/edited from callback
	respMsg := slack.NewBlockMessage(respBlocks...)

	b, err = json.Marshal(&respMsg)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))
}
