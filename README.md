JavaMacLauncher 1.1.0-SNAPSHOT
==============================

This application is a native Java application launcher written in [Go](https://golang.org/), that is intended for integration into a [macOS application bundle](https://en.wikipedia.org/wiki/Bundle_(macOS)). 

In times of Java 6, when Apple provided a Java Runtime Environment for macOS, there was a similar solution called *JavaApplicationStub*, which is unfortunately not usable with recent Java versions. Up to [macOS 10.15 (Catalina)](https://en.wikipedia.org/wiki/MacOS_Catalina) a binary was not necessary anymore to start a Java application from an application bundle. A simple Bash script like [universalJavaApplicationStub](https://github.com/tofi86/universalJavaApplicationStub) could do the job.

Apple made some changes with [macOS 11 (Big Sur)](https://en.wikipedia.org/wiki/MacOS_Big_Sur). It is indeed still possible to use a Bash script for the task, but it has some downsides regarding security. If the application needs further permissions from the operating system (e.g. access to certain user folders or access to the screen), the user needs to grant permissions to Bash or `/usr/bin/env` — not to the Java application itself. Granting those permissions to those general utilities would **undermine the system security** and is also **not intuitive** for the user. This also might become a showstopper for publishing your application to the AppStore.

The [universalJavaApplicationStub](https://github.com/tofi86/universalJavaApplicationStub) project addressed this issue by providing a compiled version of their Bash script. As they currently only provide binaries for Catalina and Big Sur we are not sure, if this approach also works with older versions of macOS — even if recent OpenJDK builds should work from [macOS 10.9 (Mavericks)](https://en.wikipedia.org/wiki/OS_X_Mavericks) upwards.

Inspired by [macstub](https://github.com/dfrugg/macstub) we've decided to build our own native Java launcher to address these issues. As we are mostly developing on Linux systems this launcher is written in [Go](https://golang.org/) for easier development and testing. Cross compilation for macOS is currently **not** possible.


How to use
----------

Assuming you have an application bundle called `MyApplication.app`, that provides its own Java Runtime Environment in this directory structure:

```
🗀 MyApplication.app
↳ 🗀 Contents
  ↳ Info.plist
  ↳ 🗀 MacOS
    ↳ JavaMacLauncher
  ↳ 🗀 Resources
    ↳ 🗀 bin
      ↳ java
    ↳ 🗀 conf
    ↳ 🗀 legal
    ↳ 🗀 lib
    ↳ 🗀 modules
    ↳ 🗀 share
      ↳ icon.icns
      ↳ splash.png
```

Things to consider:

- The *JavaMacLauncher* application (provided by this project) is placed in the `Contents/MacOS` folder with executable permission.

- The Java Runtime Environment with its files and folders is placed in `Contents/Resources`.

- The application's Java modules are stored in `Contents/Resources/modules`.

- The icon used by the application bundle is located at `Contents/Resources/share/icon.icns`.

- The splash image used during application startup is located at `Contents/Resources/share/splash.png`.

In this scenario the application bundle descriptor `Contents/Info.plist` would look like this:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>CFBundleInfoDictionaryVersion</key>
        <string>6.0</string>

        <key>CFBundleIdentifier</key>
        <string>com.mycompany.myapplication</string>

        <key>CFBundleName</key>
        <string>MyApplication</string>

        <key>CFBundleDisplayName</key>
        <string>MyApplicationName</string>

        <key>CFBundleVersion</key>
        <string>1.2.3</string>

        <key>CFBundlePackageType</key>
        <string>APPL</string>

        <key>CFBundleExecutable</key>
        <string>JavaMacLauncher</string>

        <key>CFBundleIconFile</key>
        <string>share/icon.icns</string>

        <key>LSMinimumSystemVersion</key>
        <string>10.9</string>

        <!-- Configure JavaMacLauncher. -->
        <key>JavaMacLauncher</key>
        <dict>

            <!-- Set java home directory. -->
            <key>JavaHome</key>
            <string>Contents/Resources</string>

            <!-- Set java command explicitly. -->
            <key>JavaCommand</key>
            <string>Contents/Resources/bin/java</string>

            <!-- Set further java options. -->
            <key>JavaOptions</key>
            <array>
                <string>-Dfile.encoding=UTF-8</string>
            </array>

            <!-- Set java class path entries. -->
            <key>JavaClassPath</key>
            <array>
                <string></string>
            </array>

            <!-- Set java module path entries. -->
            <key>JavaModulePath</key>
            <array>
                <string>Contents/Resources/modules</string>
            </array>

            <!-- Set application working directory. -->
            <key>WorkingDirectory</key>
            <string></string>

            <!-- Set application command. -->
            <key>ApplicationCommand</key>
            <string>-m com.mycompany.myapplication/com.mycompany.myapplication.MyApplication</string>

            <!-- Set further application command line arguments. -->
            <key>ApplicationArguments</key>
            <array>
                <string></string>
            </array>

            <!-- Set minimum reserved heap space. -->
            <key>HeapMinimum</key>
            <string>32m</string>

            <!-- Set maximum used heap space. -->
            <key>HeapMaximum</key>
            <string>512m</string>
          
            <!-- Set application splash image. -->
            <key>SplashImage</key>
            <string>Contents/Resources/share/splash.png</string>

            <!-- Set application name shown in the menu bar. -->
            <key>DockName</key>
            <dict>
                <key>default</key>
                <string>My Application</string>
                <key>de</key>
                <string>Mein Programm</string>
                <key>es</key>
                <string>Mi Programa</string>
                <key>fr</key>
                <string>Mon Programme</string>
            </dict>

            <!-- Set application icon shown in dock and menu bar. -->
            <key>DockIcon</key>
            <string>Contents/Resources/share/icon.icns</string>

            <!-- Enable global menu bar. -->
            <key>UseScreenMenuBar</key>
            <true/>

            <!-- Launch application in foreground. -->
            <key>LaunchInForeground</key>
            <true/>

        </dict>
    </dict>
</plist>
```


Configuration
-------------

*JavaMacLauncher* loads its configuration from the application bundle descriptor (`Info.plist`). All of its configurations are loaded from the `JavaMacLauncher` entry, which has to be `<dict>`:

```xml
<plist version="1.0">
    <dict>
        <!-- Typical configurations before as required by macOS. -->

        <key>JavaMacLauncher</key>
        <dict>
            <!-- Configurations for JavaMacLauncher. -->
        </dict>
    </dict>
</plist>
```


### Set Java home & Java command

You need to tell *JavaMacLauncher*, where to find a Java Runtime Environment to start your application.

```xml
<key>JavaHome</key>
<string></string>

<key>JavaCommand</key>
<string></string>
```

In case you are providing your own Java Runtime Environment within the application bundle, you should either set `JavaHome` or `JavaCommand` value.

- `JavaHome` points to the folder container the JRE / JDK.

- `JavaCommand` points to the `java` executable file.

- If `JavaHome` is set, the application sets the `JAVA_HOME` environment variable accordingly before starting the application.

- If neither `JavaHome` nor `JavaCommand` is set (or both are empty), the application tries to find the Java home folder automatically:
  
  - It looks for a `JAVA_HOME` environment variable, which is in most cases not available, if an application bundle is started regularly.

  - Otherwise, it launches `/usr/libexec/java_home` to determine a Java home.

- If a Java home was configured (or detected) and no `JavaCommand` was configured, the application assumes `bin/java` within the Java home directory to be the Java command.

- You might use relative paths for `JavaHome` and `JavaCommand`. These are converted to absolute paths based on the application bundle's absolute location.


### Set Java command line options

You might specify further options, that are passed to the Java Runtime Environment.

```xml
<key>JavaOptions</key>
<array>
    <string>-Dfile.encoding=UTF-8</string>
    <string>-Dcom.mycompany.myapplication.setting=example</string>
</array>
```


### Set Java class path or module path

If your application uses the old **class path approach**, you can add as much class path entries you like:

```xml
<key>JavaClassPath</key>
<array>
    <string>Contents/Resources/jars/MyApplication.jar</string>
    <string>Contents/Resources/more-jars/*</string>
</array>
```

If your application uses the new **module path approach**, you can add as much module path entries you like:

```xml
<key>JavaModulePath</key>
<array>
    <string>Contents/Resources/modules</string>
    <string>Contents/Resources/more-modules</string>
</array>
```

- You might use relative paths for `JavaClassPath` and `JavaModulePath` entries. These are converted to absolute paths based on the application bundle's absolute location.


### Set working directory

You might configure a certain working directory. *JavaMacLauncher* changes to this directory before starting the Java application.

```xml
<key>WorkingDirectory</key>
<string></string>
```

- If `WorkingDirectory` is not configured or empty, *JavaMacLauncher* changes to user's home directory. 

- You might use a relative path for `WorkingDirectory`. In this case the path is converted to an absolute path based on the application bundle's absolute location.


### Set Application start command

You need to tell *JavaMacLauncher* how to start your application. There are multiple possibilities:


#### Starting from class path

Just provide the class path of your application's main class:

```xml
<key>ApplicationCommand</key>
<string>com.mycompany.myapplication.MyApplication</string>
```


#### Starting from module path

Just provide the module and class path of your application's main class:

```xml
<key>ApplicationCommand</key>
<string>-m com.mycompany.myapplication/com.mycompany.myapplication.MyApplication</string>
```


#### Starting from jar file

Just provide the application's JAR file to load:

```xml
<key>ApplicationCommand</key>
<string>Contents/Resources/jars/MyApplication.jar</string>
```

- Do **not** provide the `-jar` parameter in this case. *JavaMacLauncher* will add this automatically, if `ApplicationCommand` ends with `.jar`.

- You might use a relative path for `ApplicationCommand`. If the value ends with `.jar`, the path is converted to an absolute path based on the application bundle's absolute location.


### Set Application command line arguments

You might add further command line arguments, that are passed to your application's `main(String[] args)` method:

```xml
<key>ApplicationArguments</key>
<array>
    <string>first argument</string>
    <string>second argument</string>
</array>
```


### Set minimum & maximum heap size

You might provide a minimum and maximum heap size for the Java Runtime Environment.

```xml
<key>HeapMinimum</key>
<string>32m</string>

<key>HeapMaximum</key>
<string>512m</string>
```

- *JavaMacLauncher* sets the `-Xms` / `-Xmx` parameters accordingly.

- For further tweaking, you might dismiss `HeapMinimum` / `HeapMaximum` and provide your custom settings within `JavaOptions`.


### Set application's splash image

You might specify the application's splash image:

```xml
<key>SplashImage</key>
<string>Contents/Resources/share/splash.png</string>
```

- If a `SplashImage` configuration is present, *JavaMacLauncher* will use the `-splash` parameter to start the Java Runtime Environment.

- You might use a relative path for `SplashImage`. In this case the path is converted to an absolute path based on the application bundle's absolute location.


### Set application's name in Dock & Menubar

You might specify the application's name shown in the Dock and Menubar:

```xml
<key>DockName</key>
<dict>
    <key>default</key>
    <string>My Application</string>
    <key>de</key>
    <string>Mein Programm</string>
    <key>es</key>
    <string>Mi Programa</string>
    <key>fr</key>
    <string>Mon Programme</string>
</dict>
```

- Beside the `default` name you can provide alternative names depending on the language used in the user's operating system.

- If no `default` name was set but translations in other languages are available, *JavaMacLauncher* will fallback to the global `CFBundleDisplayName` or `CFBundleName` as default value for unsupported languages.

- If a `DockName` configuration is present, *JavaMacLauncher* will use the `-Xdock:name` parameter to start the Java Runtime Environment.


### Set application's icon in Dock

You might specify the application's icon shown in the Dock:

```xml
<key>DockIcon</key>
<string>Contents/Resources/share/icon.icns</string>
```

- If a `DockIcon` configuration is present, *JavaMacLauncher* will use the `-Xdock:icon` parameter to start the Java Runtime Environment.

- You might use a relative path for `DockIcon`. In this case the path is converted to an absolute path based on the application bundle's absolute location.


### Enable global Menubar

You might enable the global Menubar for the application:

```xml
<key>UseScreenMenuBar</key>
<true/>
```

- If a `UseScreenMenuBar` configuration is present and set to `<true>`, *JavaMacLauncher* will use the `-Dapple.laf.useScreenMenuBar=true` parameter to start the Java Runtime Environment.


### Launch application in foreground

You might force application startup in foreground:

```xml
<key>LaunchInForeground</key>
<true/>
```

- If a `LaunchInForeground` configuration is present and set to `<true>`, *JavaMacLauncher* will start the Java Runtime Environment and wait until it was shutdown. Otherwise, *JavaMacLauncher* starts the Java Runtime Environment as background process and quits immediately.

- It seems to provide a slightly better user experience in macOS Big Sur to start the Java Runtime Environment in foreground. Otherwise, the tooltip on the Dock Icon shows `java` instead of the application name. 
  
- We are not sure about further advantages / disadvantages about running in background / foreground. Therefore, we leave it up to you which method to use. 


Supported operating systems
---------------------------

- Intended to be used on macOS (10.9 or newer). Releases are compiled for x86-64.
- Can be built for Linux or other Unix based systems, which might be useful for development and testing.


License
-------

This application is licensed under the terms of the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0.html). Take a look at [`LICENSE.txt`](https://github.com/OpenIndex/JavaMacLauncher/blob/develop/LICENSE.txt) for the license text.


Third party components
----------------------

- [Go](https://golang.org/) v1.16
  [(BSD License)](https://golang.org/LICENSE)
- [MacDriver](https://github.com/progrium/macdriver) v0.0.2
  [(MIT License)](https://raw.githubusercontent.com/progrium/macdriver/main/LICENSE)
- [go-plist](https://github.com/DHowett/go-plist) v0.0.0-20201203080718-1454fab16a06
  [(BSD License)](https://raw.githubusercontent.com/DHowett/go-plist/master/LICENSE)


Further information
-------------------

- [*JavaMacLauncher* at GitHub](https://github.com/OpenIndex/JavaMacLauncher)
- [Releases of *JavaMacLauncher*](https://github.com/OpenIndex/JavaMacLauncher/releases)
- [Changelog of *JavaMacLauncher*](https://github.com/OpenIndex/JavaMacLauncher/blob/develop/CHANGELOG.md)
