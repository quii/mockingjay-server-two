# mockingjay-server-two

The sequel (re-write) of [Mockingjay-server](https://github.com/quii/mockingjay-server)

## Hopes and dreams

- I am a much better engineer compared to when I wrote M1J, back in **2015**. It was my first Go project. I think I should be able to get more features but have a much simpler codebase than what MJ ended up being. 
- A lot of people (relatively speaking) saw value in MJ1, and enjoyed using it, but YAML as a config choice was pure pain. Hoping to do much better here
- Address some of the shortfalls, in particular the need to have some sense of variance between how the local server works, and running the CDC against the real one. Stuff like authentication, API keys, yada yada
- Take a very strict, ATDD approach as I described in [learn go with tests](https://quii.gitbook.io/learn-go-with-tests/testing-fundamentals/scaling-acceptance-tests)
- Leverage example mapping to give myself focus https://github.com/quii/mockingjay-server-two/wiki/Example-mapping
- Build a fast, simple UI where people can dynamically define their server and persist the configuration as they write it. Using https://htmx.org

## Things i'm optimising for

- MJ1 was _super_ fast, not just at spinning up, but leveraged concurrency to do the CDC checks, which was really nice
- The error messaging in MJ1 wasn't amazing at times, especially when the CDC failed, hope to improve this.
- MJ1 was simple becauuse it was very basic and strict, this was good but at times it meant it was unsuitable for some use-cases. Without turning it into a beast like Pact, I'd like to find a way of letting people add some flexibility to the way they define their contracts
