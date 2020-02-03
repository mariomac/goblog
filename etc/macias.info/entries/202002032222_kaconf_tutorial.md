# KAConf: Kick-Ass Configuration system for your java applications (Tutorial)

This blog post is a step-by-step guide that shows you the basic usage of
[KAConf](http://github.com/mariomac/kaconf), an Annotation-based configuration system
inspired in the wonderful [Spring Boot](http://spring.io), but simpler, lighter,
not so magic, and independent of any large framework, with no transitive
dependencies.
                                         
Its strong points are:
                                            
* Easy to use, integrate and extend
* Tiny footprint: a single, <13KB JAR with no third-party dependencies
* Born from own necessity, with no over-engineered use cases

## Simplest use case: Hello, world!

Imagine you want to create a `Hello, world!` application featuring high
configurability to be able to decide both the greeting (`Hello`, `Hi`, `Heyya`...)
and the receiver of the greeting (`world`, `people`, `Peter`, `friends`...).

You may first want to isolate all the configuration as a Java Bean, with
hardcoded configuration values:

```java
public class Config {
    private String greet = "Hello";
    private String name = "world";

    public final String getGreet() {
        return greet;
    }

    public final String getName() {
        return name;
    }
}
```

Then you can inject this configuration into a `Greeter` object that
makes use of it for the composition of greetings:

```java
public class Greeter {
    private final Config cfg;

    public Greeter(Config cfg) {
        this.cfg = cfg;
    }

    public String greet() {
        return cfg.getGreet() + ", " + cfg.getName() + "!";
    }
}
```

And then you instantiate and relate `Config` and `Greeter`, and
make use of your `Greeter` instance:

```java
public class Main {
    public static void main(String[] args) {
        var cfg = new Config();
        var greeter = new Greeter(cfg);

        System.out.println(greeter.greet());
    }
}
```

When you run your application (e.g. with the Gradle application plugin), you will finally see a greeting
message in your terminal:

```
$ ./gradlew run 
Hello, world!
```

## Adding KAConf support

As hardcoded configurations are not really useful but for default values, you may want to
be able to provide your own external configuration via a properties' file such as the
following example:

```properties
greeter.greet = Wassssup
greeter.name = my super friend
```

KAConf allows you to just map those properties to your `Config` bean by using simple
annotations.

First, you need to add the `info.macias:kaconf:0.9.0` (the newest version at the moment
of writing this) to your Gradle or Maven project:

Gradle (Kotlin):
```kotlin
dependencies {
    compile("info.macias", "kaconf", "0.9.0")
}
```

Maven:

```XML
<dependency>
    <groupId>info.macias</groupId>
    <artifactId>kaconf</artifactId>
    <version>0.9.0</version>
</dependency>
```

Then, you are now able to annotate your Java Bean properties with the
`@Property` annotation:

```java
import info.macias.kaconf.Property;

public class Config {
    @Property("greeter.greet")
    private String greet = "Hello";

    @Property("greeter.name")
    private String name = "world";

    public final String getGreet() {
        return greet;
    }

    public final String getName() {
        return name;
    }
}
```
(You may want to keep the hardcoded `"Hello"` and `"world"` as default values
when any configuration property is missing)

To allow KAConf mapping your properties to your Java Bean, you need to create
a `Configurator` object via a `ConfiguratorBuilder` instance that receives a
`Properties` object loaded from the file passed as argument, and then
invoke the `Configurator.configure(...)` method passing there the `Config`
object:

```java
import java.io.FileInputStream;
import java.util.Properties;
import info.macias.kaconf.ConfiguratorBuilder;

public class Main {
    public static void main(String[] args) throws Exception {
        var cfg = new Config();

        if (args.length > 0) {
            var cfgBuilder = new ConfiguratorBuilder();
            var p = new Properties();
            p.load(new FileInputStream(args[0]));
            cfgBuilder.addSource(p);
            cfgBuilder.build().configure(cfg);
        }

        var greeter = new Greeter(cfg);
        System.out.println(greeter.greet());
    }
}
```

KAConf will look for the `@Property` annotations in the `Config` object and
replace there the matching properties, if any. Then you can see how your
application is configured at runtime:

```
$ ./gradlew run --args="config.properties"
Wassssup, my super friend!
```

## Multiple configuration sources

You may want to configure your application from different,
prioritized sources; e.g. environment variables (the highest priority), configuration
file, and default, hardcoded values (the lowest priority).

In addition, different configuration sources may follow different naming
conventions (e.g. the `greeter.greet` property from the Java properties' file
would usually be named `GREETER_GREET` if it is defined as an environment variable).

For that reason, you can specify multiple names in the `Property` annotation:

```java
import info.macias.kaconf.Property;

public class Config {
    @Property({"GREETER_GREET", "greeter.greet"})
    private String greet = "Hello";

    @Property({"GREETER_NAME", "greeter.name"})
    private String name = "world";

    public final String getGreet() {
        return greet;
    }

    public final String getName() {
        return name;
    }
}
```

And you can add to the `ConfiguratorBuilder` multiple sources of configuration
in order of priority (highest first):

```java
import java.io.FileInputStream;
import java.util.Properties;
import info.macias.kaconf.ConfiguratorBuilder;

public class Main {
    public static void main(String[] args) throws Exception {
        var cfg = new Config();

        var cfgBuilder = new ConfiguratorBuilder()
                .addSource(System.getenv());

        if (args.length > 0) {
            var p = new Properties();
            p.load(new FileInputStream(args[0]));
            cfgBuilder.addSource(p);
        }
        cfgBuilder.build().configure(cfg);

        var greeter = new Greeter(cfg);
        System.out.println(greeter.greet());
    }
}
```

You can see how your application now is configured according to the given
priorities:

```
$ export GREETER_GREET="Hey"
$ ./gradlew clean run --args="config.properties"
Hey, my super friend!
$ export GREETER_NAME="you"
$ ./gradlew clean run --args="config.properties"
Hey, you!
$ unset GREETER_GREET
$ ./gradlew clean run --args="config.properties"
Wassssup, you!
```


