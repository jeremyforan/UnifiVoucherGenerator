
# UniFi Network Application - Go Voucher Generator

**Origin Story**: After setting up Ubiquiti access points at my place, I aimed to offer straightforward Guest Wi-Fi. I didn't want to use a PSK. However, I also needed to ensure that only people inside my home could use it, avoiding external access. To solve this, I introduced a button that generates a small voucher code upon being pressed.

The is an open-source Golang library designed to interface with the UniFi Network Application dashboard, facilitating the automated generation and fetching of new vouchers for WiFi hotspot landing pages.

By enabling the programmatic generation of vouchers, this library broadens the scope for creating better guest Wi-Fi experiences.

![](demo.gif)

## Usage/Examples
A basic example of the library creating a basic voucher.

```go
package main

import (
	"fmt"
	"github.com/jeremyforan/UnifiVoucherGenerator"
	"github.com/jeremyforan/UnifiVoucherGenerator/voucher"
	"log/slog"
	"net/url"
)

func main() {

	// URL of the Unifi controller.
	baseUrl, err := url.Parse("https://127.0.0.1:8443")
	if err != nil {
		panic(err)
	}

	// Create a new client with the username, password, and URL of the Unifi controller.
	client := UnifiVoucherGenerator.NewClient("user@email.com", "p455w0rd", baseUrl)

	// Login to the Unifi controller.
	err = client.Login()
	if err != nil {
		slog.Info("Failed to login to Unifi controller")
		panic(err)
	}

	// Create a new voucher with default settings.
	v := voucher.NewDefaultVoucher()

	// Add the Voucher to the Unifi controller.
	err = client.AddVoucher(v)
	if err != nil {
		slog.Error("unable to add voucher")
	}

	fmt.Println(v.AccessCode())
}
```


## Documentation

### Expiry

### Network Limits

### Usage



## Running Tests

To run tests, run the following command

```bash
  npm run test
```


## Deployment

To deploy this project run

```bash
  npm run deploy
```

## Limitations


## Acknowledgements

- [Awesome Readme Templates](https://awesomeopensource.com/project/elangosundar/awesome-README-templates)
- [Awesome README](https://github.com/matiassingers/awesome-readme)
- [How to write a Good readme](https://bulldogjob.com/news/449-how-to-write-a-good-readme-for-your-github-project)


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`API_KEY`

`ANOTHER_API_KEY`


## FAQ

#### Question 1

Answer 1

#### Question 2

Answer 2


## Feedback

If you have any feedback, please reach out to us at jeremy.foran@gmail.com


## License

[MIT](https://choosealicense.com/licenses/mit/)

