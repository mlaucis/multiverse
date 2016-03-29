# Building iOS frameworks with storyboards, nibs and resources

For the last month we have been working on creating [Tapglue Elements](https://github.com/tapglue/elements-ios), a framework on iOS for integrating full features into third party apps. This includes GUIs and graphical assets like images and also interacting with our current SDK which is responsible for networking and caching. One of the fundamental requirements was the support of the biggest dependency management tool out there for iOS: CocoaPods.

Too our big surprise there were not many examples of frameworks that do something similar to our goal, create a full feature with its own designs and UX to be integrated into third party apps. The most similar examples we found actually avoided using storyboards and xib files all together and did views purely in code. 

In this post our goal is to create a step-by-step guide on how to create a framework with storyboards, xibs, assets and localisation that works with CocoaPods. A project containing all the code can be found [here](https://github.com/nilsen340/ios-framework-with-storyboard).

## Setup

To start we recommend using CocoaPods own template for creating frameworks by running

`pod lib create MyFramework`

You can find further documentation [here](https://guides.cocoapods.org/making/using-pod-lib-create.html).

We will be using Swift and a demo app, for the sake of simplicity, will not use any test frameworks.

This sets up a framework connected to a demo application you can then use to try out our framework with. In our case it looks like this:

![alt tag](http://i.imgsafe.org/8b4e5cb.png)

This also generates the podspec for us, for now we'll just leave it the way it is. 

## Lets get started!

Lets start by creating a storyboard and a view controller with a tableview. Embed your view controller in a `UINavigationViewController` and make the navigation view controller the initial view controller.

![alt tag](http://i.imgsafe.org/f0c0d5d.png)

We link this with the view controller we will create  called `FrameworkVC`. To spice things up we will implement our cells in xibs, we name it `OurCell.xib` and Assign it the reuse identifier `OurCell`

![alt tag](http://i.imgsafe.org/2b0f60e.png)

Now lets have a look at how we can combine all of this together in the view controller:

```swift
import UIKit

public class FrameworkVC: UIViewController {

    @IBOutlet weak var tableView: UITableView!
    override public func viewDidLoad() {
        super.viewDidLoad()
        let podBundle = NSBundle(forClass: FrameworkVC.self)
        
        let bundleURL = podBundle.URLForResource("MyFramework", withExtension: "bundle")
        let bundle = NSBundle(URL: bundleURL!)!
        let cellNib = UINib(nibName: "OurCell", bundle: bundle)
        tableView.registerNib(cellNib, forCellReuseIdentifier: "OurCell")
        tableView.estimatedRowHeight = 80
        tableView.rowHeight = UITableViewAutomaticDimension
    }
}

extension FrameworkVC: UITableViewDelegate {}

extension FrameworkVC: UITableViewDataSource {
    public func tableView(tableView: UITableView, cellForRowAtIndexPath indexPath: NSIndexPath) -> UITableViewCell {
        return tableView.dequeueReusableCellWithIdentifier("OurCell")!
    }
    
    public func tableView(tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return 1
    }
}
```

To load the xib file we need to use a bundle associated to the framework itself, hence in the `viewDidLoad` we ask for the bundle of the FrameworkVC and get the bundle of the same name as we declared in the podspec. At this point we need to head over to the podspec to update it to include our storyboard and xib file.

```ruby
Pod::Spec.new do |s|
  s.name             = "MyFramework"
  s.version          = "0.1.0"
  s.summary          = "A short description of MyFramework."
  s.description      = <<-DESC
                       DESC

  s.homepage         = "https://github.com/<GITHUB_USERNAME>/MyFramework"
  # s.screenshots     = "www.example.com/screenshots_1", "www.example.com/screenshots_2"
  s.license          = 'MIT'
  s.author           = { "John Nilsen" => "john@tapglue.com" }
  s.source           = { :git => "https://github.com/<GITHUB_USERNAME>/MyFramework.git", :tag => s.version.to_s }
  # s.social_media_url = 'https://twitter.com/<TWITTER_USERNAME>'

  s.platform     = :ios, '8.0'
  s.requires_arc = true

  s.source_files = 'Pod/Classes/**/*.{swift}'
  s.resource_bundles = {
    'MyFramework' => ['Pod/Classes/**/*.{storyboard,xib}']
  }
end
```

There are two significant changes here, on one side we filter the source files by extension by adding `.{swift}` and since we're putting storyboards and nibs in the classes folder we changed the `resource_bundle` to `Pod/Classes/**/*.{storyboard,xib}`

Now we head over to the Example app and executes `pod install`. In Xcode the storyboards and xibs will now be displayed in a different group: `Resources`

## Hooking the pod into the app

Since this is a framework we need to hook it up into the demo app to be able to see the results of our work. Lets start by adding a segue into our view controller. For this I will add a new class: `MyFramework` where we can add the segue call.

```swift
import UIKit

public class MyFramework {
    
    public static func performSegueToFrameworkVC(caller: UIViewController) {
        let podBundle = NSBundle(forClass: FrameworkVC.self)
        
        let bundleURL = podBundle.URLForResource("MyFramework", withExtension: "bundle")
        let bundle = NSBundle(URL: bundleURL!)!
        let storyboard = UIStoryboard(name: "FrameworkStoryboard", bundle: bundle)
        let vc = storyboard.instantiateInitialViewController()!
        caller.presentViewController(vc, animated: true, completion: nil)
    }
}
```

As you can see we wrote the exact same code for getting the bundle, to tidy up we will create a method to generate it for us inside the `MyFramework` class.

```swift
    static var bundle:NSBundle {
        let podBundle = NSBundle(forClass: FrameworkVC.self)
        
        let bundleURL = podBundle.URLForResource("MyFramework", withExtension: "bundle")
        return NSBundle(URL: bundleURL!)!
    }
```

With that refactoring we can tidy up both `MyFramework` and `FrameworkVC`.

Thats it! Now we can hook it up into the demo app. If we head over to the `ViewController` in the `Example for MyFramework` folder we just need to add the segue call in `viewDidAppear`

```swift
import UIKit
import MyFramework

class ViewController: UIViewController {
    
    override func viewDidAppear(animated: Bool) {
        MyFramework.performSegueToFrameworkVC(self)
    }
}
```

When running ti you should see something like this:

![alt tag](http://i.imgsafe.org/e1b361c.png)

## Adding images into the mix

Lets improve our cell design by adding an image to it. First we create an asset catalog in the Classes folder of MyFramework. Then we press plus to add a image set. Then we add the following images.

[download images](https://github.com/nilsen340/ios-framework-with-storyboard/raw/master/tapglue-logo.zip)

Lets redesign our cell to look like this

![alt tag](http://i.imgsafe.org/f6dd016.png)

If we were to run pod install on our project it would break. We need to add these new files to the resource bundle.

Our podspec should look something like this after fixing the issue:

```ruby
Pod::Spec.new do |s|
  s.name             = "MyFramework"
  s.version          = "0.1.0"
  s.summary          = "A short description of MyFramework."
  s.description      = <<-DESC
                       DESC

  s.homepage         = "https://github.com/<GITHUB_USERNAME>/MyFramework"
  # s.screenshots     = "www.example.com/screenshots_1", "www.example.com/screenshots_2"
  s.license          = 'MIT'
  s.author           = { "John Nilsen" => "john@tapglue.com" }
  s.source           = { :git => "https://github.com/<GITHUB_USERNAME>/MyFramework.git", :tag => s.version.to_s }
  # s.social_media_url = 'https://twitter.com/<TWITTER_USERNAME>'

  s.platform     = :ios, '8.0'
  s.requires_arc = true

  s.source_files = 'Pod/Classes/**/*.{swift}'
  s.resource_bundles = {
    'MyFramework' => ['Pod/Classes/**/*.{storyboard,xib,xcassets,json,imageset,png}']
  }
end
```

Notice the `resource_bundle` now includes extensions like `xcassets`, `json`, `imageset` and `png` in addition to the ones from earlier.

Now run `pod install` from the Example folder and we're all set again!

## Tips

Some of the minor issues I ran into were related to the pod not being updated when executing `pod spec lint`, we usually solved all of these doing a `pod cache clean --all` and executing `pod spec lint` again. 

When writing a framework like this we would recommend providing the view controllers themselves and not the segues, and provide delegation of the most relevant parts of the view controllers. Thats the approach we decided to use for [Tapglue elements](https://github.com/tapglue/Elements-ios)

## Wrapping up

Thats it! If you want further examples of how to implement this I recommend you read the [Tapglue Elements](https://github.com/tapglue/elements-ios) source code. 
