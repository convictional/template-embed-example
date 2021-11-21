# Using go:embed to get compile-time guarantees for templates
Most go projects that need to either send emails or display web pages will be using Go's native templating language to do so. Go has a very powerful templating language that has made dealing with templating data in HTML extremely convenient.

You will be able to find the code for this post [here](https://github.com/convictional/template-embed-example) and each section will contain a git tag for the specific part of the codebase at that point

## Standard go template example
##### ef995ebdbcd89015fd8803cf2a4997862592e7f4
First let us put together a standard usage of templating in go. This is what most projects that use templating are going to look at, and what our first implementation of emailing looked like:

```go  
// Read in template (in this example we are sending a forgot password email)
passwordTemplate, err := template.ParseFiles(fmt.Sprintf("%s/%s", basepath, templateLayout), fmt.Sprintf("%s/%s", basepath, templateForgotPassword))
if err != nil {
    panic(err)
}
// Execute template with data and store in a bytes.Buffer for use in email
var body bytes.Buffer
err = passwordTemplate.ExecuteTemplate(&body, "layout", &ForgotPasswordData{Link: "https://httpbin.org"})
if err != nil {
    panic(err)
}
//[...] Send email with the body that was generated from passwordTemplate.ExecuteTemplate  
```  

This works well and will correctly template in the values we need - but there are a few drawbacks with using this method. First - the template we are loading in is parsed from the file system each time we want to send an email, which is very inefficient if you are sending this email more often than the application restarts. Second - because the email is only sent at runtime, we have no way to guarantee that the template is there. If we are running this application in a docker container, and we incorrectly copy over the files, this email will fail at runtime and may cause issue for our customers. Finally - in order to use this email package anywhere in the module we need to use some hacky path manipulation to allow for both local testing and 'external' usage of the email package.

## Optimizing our template usage
##### f7ae625607fff46a517511b3a02b2dd522d03a8b
We can solve the first problem easily by utilizing a lesser-know feature of Go:  `init()`. `init()` is a function that automatically runs when a package is imported for the first time. What we can do with this function is initialize all of our templates when the email package is imported for the first time, and simply execute those packages when we need to put data in them.

First we need to define the template at the top of the email.go file
```go
//...
var passwordTemplateInit *template.Template
//...
```

Next we create an `init()` file that parses the template files and assigns that to the variable we just created
```go
//...  
func init() {
    passwordTemplateInit = template.Must(template.ParseFiles(fmt.Sprintf("%s/%s", basepath, templateLayout), fmt.Sprintf("%s/%s", basepath, templateForgotPassword)))
}
//...  
```  

And finally we will modify our `SendForgotPasswordEmail` function to use this new global variable rather than doing the parsing on its own

```go
//...  
// Here I created a new function rather than modifying the old one so that we can compare visually and benchmark them together  
func (s Sender) SendForgotPasswordEmailInit(address string) error {
    // Execute template with data and store in a bytes.Buffer for use in email
    var body bytes.Buffer
    err := passwordTemplateInit.ExecuteTemplate(&body, "layout", &ForgotPasswordData{Link: "https://httpbin.org"})
    if err != nil {
        panic(err)
    }
    return s.sendEmail(address, "Reset Password", body.String())
} 
//...  
```  

You can immediately see that the function is much simpler now that we have moved the template initialization code out of the function that actually sends the email. This is also **exceptionally** faster:

```sh  
# Lower ns/op is better  
go test -bench=.  
BenchmarkSendForgotPasswordEmailTemplatesInit-16    	  507255	      2139 ns/op
BenchmarkSendForgotPasswordEmailTemplates-16        	   12099	    101731 ns/op
PASS
ok  	github.com/convictional/template-embed-example/email	3.524s 
```  

There are still a few problems with this approach however, since we are still vulnerable to a copying mistake or a deletion of a template file causing runtime errors, and we do still have to work around where in the filesystem the templates are located.

## Introducing go:embed
##### b3b1e0dfe1e6e38e6ce5e5b6e952f85d881d7311
The `go:embed` directive was introduced in version 1.16 alongside the [virtual FS proposal](https://go.googlesource.com/proposal/+/master/design/draft-iofs.md_) that was also introduced in 1.16. By attaching an absolute route to a variable with the `go:embed` directive, the go compiler will automatically load in the files located at the route into the defined variable. If a file is not found at that location, the compiler will throw an error. Since this is an absolute path, we can also remove the hacky pathing solutions we needed in order to use the email package elsewhere in the module, simplifying our code a lot.

The first step to making this change is to declare some new variables:
```go
//...  
var (
    //go:embed templates/layout.html
    baseLayoutFS embed.FS
    //go:embed templates/forgot_password.html
    passwordTemplateFS embed.FS
    passwordTemplateInitFS *template.Template
)
//...
```  
> Note the `//go:embed` directive above `baseLayoutFS` and `passwordTemplateFS`. In this case rather than embedding the file to a string or []byte we are using the `embed.FS` type. The reason we are doing this is the `html/template` package has [native support](https://pkg.go.dev/html/template#ParseFS) for parsing `embed.FS` objects.

Next we need to modify our `init()` function to handle these new variables and  ensure that they are properly set up for use
```go  
func init() {  
    // We separate baseLayout from passwordTemplateInitFS due to how ParseFS works - this is a better pattern regardless
    baseLayout := template.Must(template.New("layout").ParseFS(baseLayoutFS, templateLayout))
    passwordTemplateInitFS = template.Must(baseLayout.ParseFS(passwordTemplateFS, templateForgotPassword))
}  
```  
The rest of `email.go` only needs changes to use the new variables, and we now have  compile-time guarantees that the files we need for email templating exist and are ready for use.

Running a quick benchmark shows that we haven't lost any performance with this change:

```shell  
go test -bench=.
BenchmarkSendForgotPasswordEmailTemplatesInitFS-16    	  521109	      2105 ns/op
BenchmarkSendForgotPasswordEmailTemplatesInit-16      	  475965	      2208 ns/op
BenchmarkSendForgotPasswordEmailTemplates-16          	   10000	    102306 ns/op
PASS
ok  	github.com/convictional/template-embed-example/email	3.385s
```  

## Production Readiness
The code used in this post is simplified to make comparisons easier and to illustrate the changes themselves, not any specific structure. In the interest of utility here are some changes that would make this even more robust and easy to use in a production environment.

The first step would be to turn the email package into a library package rather than a domain package. This is somewhat of a stylistic point, but the email package should be handling sending emails, and the domain packages should be handling domain-specific wording and data templating.

To do this, we will move the forgot_password.html template (and it's init code) to the `user` package and create a function `MustParseContentFS` in the email package. `MustParseContentFS` takes in a template from an external package and combines it with the `layout.html` that is still within the `email` package.

Another change is to move the `init()` functions to an `init.go` file in each package. This isn't needed (Go will run it first regardless) but it does make it more obvious that there is some pre-work being done in the package.

You can see what this would look like (plus a few other small tweaks) on a [separate branch](https://github.com/convictional/template-embed-example/tree/production-readiness) in the repo for this post.

## Conclusions
And with that we have successfully refactored our template usage to have compile-time guarantees that our templates exist, simplified access from across the project, and optimized the parsing and execution by an order of magnitude.