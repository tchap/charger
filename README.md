# Charger

Charge your Go (Golang) object with new values full of energy!

## Motivation

I needed to fill a configuration struct from environment variables,
but no existing library was really able to do what I needed.

In particular, all libraries that I came across let you specify default values
statically. This approach breaks down when using Docker containers, because you
can set default values either in your executable or in the container when you
are building it.

That would be all good, because who cares really. When run, the executable
would simply pick the values as set by the container. Well, not so simply.
What I needed was to implement a help command for my executable placed as an
entry point in a container that would list all significant environment
variables with their default values. In this scenario it is not enough to have
some values compiled in, we need to take in the container defaults somehow.

## Ideas

- [ ] Separate charging into getting/rewriting/rendering values and make sure
      users can hook into the process at any point.
- [ ] Hook into process by either registering a global fuction for the phase
      or implement particular interface by the key spec object.
- [ ] Must support charging objects without tags.
- [ ] Value templates that can refer to other values (through interface?).
- [ ] Hyerarchical key specification, e.g. `MQTT_` prefix.
- [ ] Method for dumping config help into an io.Writer.
