## Code Quality
- "Anonymize": Remove event fields that aren't necessary, like location

## Devops
- Figure out how to make into a static binary? or just dockerize it :)


## Open-Source Contribution Opportunity

ICalendar events with multiline properties **need** to be unfolded. 
See if you can help document unfolded parsing error or fail with better message -> https://github.com/hoodie/icalendar-rs/issues/47

I could help by adding an example in the parsing docs. Now that I've spent a few hours on this, I realize that there's an examples/
dir in the repo.
Being pretty new to Rust, I didn't know to go looking in the repo's examples/ dir, since docs.rs tends to be all I need for more mature crates.
It also took me a while to find the repo itself, as the [docs.rs](https://docs.rs/icalendar/latest/icalendar/index.html) page 
doesn't seem to have links to the github

The below bit *does not* work:

```rust
let calendar_str = fs::read_to_string("examples/error_calendar.txt").unwrap();
let calendar =
    icalendar::parser::read_calendar(&calendar_str).unwrap_or_else(|err| panic!("{}", err));
```

The below bit works:

```rust
let calendar_str = fs::read_to_string("examples/error_calendar.txt").unwrap();
let unfolded = &icalendar::parser::unfold(&calendar_str);
let calendar =
    icalendar::parser::read_calendar(&unfolded).unwrap_or_else(|err| panic!("{}", err));
```
