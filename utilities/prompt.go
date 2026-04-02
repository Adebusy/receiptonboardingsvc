package utilities

import "fmt"

func NewPromptForReceipt(jsonData string, logoURL string) string {
	const systemPrompt = `
You are a professional Invoice Generator. You will receive a JSON object containing transaction details. 
Your goal is to return a single, valid, and responsive HTML string for a retail receipt.

### Design Specifications:
1. **Header**: Center-aligned. Include the company logo at the top using <img src='%s' style='width:120px;'>.
2. **Typography**: Use a clean sans-serif font for the body. Use a bold font for 'ADEBUSY CONCEPT'.
3. **Handwriting Effect**: For the 'Amount in Words' section, use the Google Font 'Dancing Script'. 
   Include this link in the HTML head: <link href='https://fonts.googleapis.com/css2?family=Dancing+Script&display=swap' rel='stylesheet'>.
4. **Table**: Create a clean table with a navy blue header (#002D62) and white text. Add alternating light gray background rows for readability.
5. **Totals**: Display the Grand Total in a large, bold font at the bottom right.

### Business Data:
- **Company**: Adebusy Concept, 25, Repton Drive, Coventry.
- **Manager**: Alao Ramon.
- **Transaction Data**: %s

### Constraint:
Return ONLY the raw HTML code. Do NOT include any markdown formatting like ` + "```" + `html or conversational text. 
The output must be ready to be saved as an .html file immediately.`

	return fmt.Sprintf(systemPrompt, logoURL, jsonData)
}

func LogoConceptPrompt(companyName, specialization string) string {
	return fmt.Sprintf(`
You are a senior brand designer and creative director with strong experience in corporate identity and visual branding.

Your task is to creatively design and describe a minimum of THREE (3) distinct company logo concepts based on the following information:

Company name: %s
Industry / Area of specialisation: %s

Design objectives:
- Produce modern, clean, and professional logo concepts
- Logos must be suitable for corporate, legal, and commercial registration
- Designs should scale well across digital platforms and print materials
- Each concept must be clearly different in visual style, symbolism, structure, and color usage
- Avoid generic or overused design clichés
- Prioritise originality and visual memorability while maintaining brand credibility

For EACH logo concept, provide the following fields:
- "logo_id" (unique identifier)
- "concept_name" %s
- "visual_description" (detailed explanation of the logo’s visual idea and symbolism)
- "color_palette" (array of primary brand colors in HEX format)
- "logo_type" (choose one: wordmark, symbol, abstract, combination mark)
- "brand_personality" (array describing emotional and brand attributes)
- "logo_base64" object containing:
  - "extension" (e.g. svg)
  - "data" (a valid Base64-encoded representation of the suggested logo)

Output requirements:
- Return the response strictly as a valid JSON array
- Include at least THREE (3) logo concept objects
- Do NOT include explanations, comments, or text outside the JSON structure
- Ensure Base64 strings are valid and decodable

Creativity is strongly encouraged.
`, companyName, specialization, companyName)
}

func CompanyNamesPrompt(companyName, specialization string) string {
	return fmt.Sprintf(`
Suggest the top 3 professional company names based on the following details:

Founder name: %s
Area of specialisation: %s

The names should be:
- Modern and professional
- Easy to pronounce
- Suitable for a registered company
- Unique and brandable

Return only the company names as a list.
`, companyName, specialization)
}

func ReceiptPrompt(customer_name, total_price, amount_in_word string, items_list string) string {
	retJsonItem, TotalPrice := ParseItems(items_list)

	return fmt.Sprintf(`You are a professional document architect. Your task is to generate a single-file, responsive HTML/CSS invoice.

		Requirements:

		Use modern, inline CSS only (no external files).

		Primary Color: Navy Blue (#002D62).

		Use a clean sans-serif font for the body and a cursive-style Google Font (like 'Dancing Script') for the 'Amount in Words' section to mimic handwriting.

		Insert a placeholder img tag with id='logo' for the company logo.

		Make all headers bold and the layout professional for a retail tech business.

		Data to include:

		Company: Adebusy Concept, 25 Repton Drive, Coventry.

		Manager: Alao Ramon.

		Customer: %s.

		Items: %s.

		Total: %s.

		Amount in Word: %s (in handwriting style).

		Return ONLY the HTML code. No conversational text.`, customer_name, retJsonItem, TotalPrice, amount_in_word)
	// {
	// 	"customer": "Splendor Concept",
	// 	"items": [{"name": "Macbook Pro", "price": 1000}, {"name": "Mouse", "price": 200}]
	//   }

	//github.com/nleeper/gowords (number to world)
}

func CompanySignaturePrompt(fullName string) string {
	return fmt.Sprintf(`Generate three distinct, realistic handwritten signature samples based on the full name "Alao Ramon".

Style requirements (apply to all signatures):
- Black ink
- Natural pen pressure with visible stroke variation
- Slight rightward slant
- Professional and executive tone
- Transparent background
- High resolution
- Output image format: PNG

Return the response strictly in valid JSON format.
Do not include any explanatory text outside the JSON.

Each signature image must be:
- Encoded as a Base64 string
- Represent a PNG image with transparency
- Visually distinct in style

Use the following JSON structure exactly:

{
  "full_name": "%s",
  "signatures": [
    {
      "style": "Executive Classic",
      "description": "Clean, balanced cursive with confident flow",
      "image_format": "png",
      "image_base64": "<BASE64_STRING>"
    },
    {
      "style": "Modern Professional",
      "description": "Minimalist, smooth curves, digital-friendly",
      "image_format": "png",
      "image_base64": "<BASE64_STRING>"
    },
    {
      "style": "Bold Authority",
      "description": "Stronger strokes, pronounced initials, authoritative finish",
      "image_format": "png",
      "image_base64": "<BASE64_STRING>"
    }
  ]
}
`, fullName)

}
