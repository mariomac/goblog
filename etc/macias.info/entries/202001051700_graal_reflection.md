# Enabling Java Reflection on GraalVM Ahead-of-Time compiler

In my [previous blog post](./201912201300_graal_aot.md) I evaluated the
feasibility of Java for lightweight system programming thanks to the
[GraalVM](https://www.graalvm.org/)
native image generation tool. Despite the initial results look promising,
I felt disappointed when they
pointed me out that GraalVM does not support reflection by default (which IMHO
is a wonderful and powerful tool to enhance the expressiveness of our software
and to reduce boilerplate). However, you can actually configure the
ahead-of-time compiler to incorporate a user-provided reflection metadata.

Let's see how.

## Simple reflection example: read KickAss Configuration data

As a quick example, let's write a simple program that can read external
configuration data using [KAConf: KickAss Configuration](https://github.com/mariomac/kaconf),
a library made by myself to easily configure software through Java Annotations
(inspired in the Spring Framework). With KAConf, you can configure your data
annotating your fields with a `@Property` annotation:

```java
import info.macias.kaconf.Property;

public class Config {
   @Property("NAME")
   private String name = "world";

   public String getName() {
  	return name;
   }
}
```

You can also annotate static fields:

```java
public class StaticConfig {
   @Property("GREETING")
   public static String GREETING = "Hello";
}
```

Then, you only have to create a `Configurator` object, tell it about the sources
of configuration data (e.g. environment) and which objects or classes are
configurable:

```java
public static void main(String[] args) {
   Configurator c = new ConfiguratorBuilder()
     	.addSource(System.getenv())
     	.build();

   Config cfg = new Config();
   c.configure(cfg);
   c.configure(StaticConfig.class);

   System.out.println(StaticConfig.GREETING + ", " + cfg.getName() + "!");
}
```

If you compile and run the above code as a normal JVM bytecode (e.g. a JAR), it
will make use of the default configuration:

```
$ java -jar aot-reflection-test.jar
Hello, world!
```

Or will read the configuration from the environment:

```
$ GREETING="What's up" NAME="gang" java -jar aot-reflection-test.jar
What's up, gang!
```

The [KAConf library](http://github.com/mariomac/kaconf) makes intensive use of
Java Reflection to automatically set the configurable values to their respective
fields. However, if you use the `native-image` GraalVM command, you can run the
program but the generated binary is not able to modify the `Config` and
`StaticConfig` values at runtime, showing always the default _Hello, world!_ 
message.

## How to pass the reflection configuration to `native-image`

First, you need to create a JSON file (e.g. `reflectconfig.json`) where you
specify the classes and fields reflection information:

```json
[
  {
	"name" : "info.macias.aotrt.StaticConfig",
	"allDeclaredFields": true
  },
  {
	"name" : "info.macias.aotrt.Config",
	"allDeclaredFields": true
  }
]
```

For the sake of simplicity, we have used the `allDeclaredFields` property as a
wildcard, but you can actually specify which concrete fields, constructors and
methods can be used for reflection ([see the online GraalVM documentation for more details](https://github.com/oracle/graal/blob/master/substratevm/REFLECTION.md)).

Then, you have to compile the bytecode into native by passing the
`-H:ReflectionConfigurationFiles` parameter to the `native-image` command:

```
$ native-image --no-fallback -H:ReflectionConfigurationFiles=reflectconfig.json\
    	-jar build/libs/aot-reflection-test.jar aot-test
```

And, _voilà_! your native application can make use of the coolest configuration
library around the world.

```
$ GREETING="What's up" NAME="gang" ./aot-test
What's up, gang!
```

## Known limitations
The ideal way of using [KAConf](https://github.com/mariomac/kaconf) is along
with immutable properties, e.g.:

```java
@Property("NAME")
private final String name = "world";

@Property("GREETING")
public static final String GREETING = KA.def("Hello");
```

However, the GraalVM `native-image` AoT compiler (at least the version
1.0.0 RC16), doesn't seem to like final fields, even if you enable writing with
the `"allowWriting": true` property inside the `reflectconfig.json` file.

## Example source code

[https://github.com/mariomac/aot-reflection-test](https://github.com/mariomac/aot-reflection-test)
