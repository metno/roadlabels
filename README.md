# Road Annotation App

This is an app to annotate images from road side web cameras with a pre defined set of classes. 
 

## Table of contents

- Objective
- How to annotate
- Classes 
- Examples
- Bugs And Restrictions 

## objective

The objective is to develop an automatic tool for detecting the current road conditions from roadside web camera images, using machine learning. To get a sufficient training dataset, we need to manually annotate a large amount of images with a set of road condition classes (1). We want samples from all seasons, light conditions and classes and on as may different cameras as time allows. Ideally we want to try to have a "balanced" dataset - About the same amount of samples from each class. At least hunt for images from the minority classes .You will se a count on the front page . For exampler "1 Dry" is plenty but "5 Heavy slush" is sparse

## How to annotate

Go to https://modellprod.k8s.met.no/roadlabels/login/ . Create an account and log in
- 1. Chose images to annotate: To do that you can Click Camera Listing . Chose a camera semi-randomly. For example where you live, on the way to the cabin, childhood memories or simplly a beautiful road . Prefer cameras with few labhels . See the annotation counts in the camera listing.  
- 2. When you have chosen a camera click on the thumb and annotate image with road condition. Click Right arrow to save. Repeat until bored or samples are very simliar. Then select another camera. Or another season/month/day with the date picker.
- 3. Ice . Ice is special and hard to find, so the annotation app has a link "Random list of images from stations with ice in frost" to help identify pictures with ice. Use that one for ice . Ice pictures are from the period 2023-02-10 - 2023-05-12 . Pr. 2023-08-31, 30 SVV stations has road_ice_thickness . Note that if Statens Vegvesens sensors indicates ice, the sensor is not necessarily correct.  

      If an image has both snow and ice, chose ice.

      If the road has been plowed but road still not visible, use class 8 Heavy snow.

- 4. When in doubt it is allowed to move the next and previous image (+/- 6hours) and check the road and its surroundings.   

### The classes (1): 
0 - Dry (No visible water or patches on the surface. Surrounding environment also looks dryish)

1 - Patchy water (Patches of water or not dry looking, with moistness in the surface)

2 - Light water (Thinner film of water, light reflection, often wet or snowy surroundings)

3 - Heavy water (Clearly very wet)

4 - Ice (Very difficult)

5 - Light slush (Small amounts of snow on a wet surface)

6 - Heavy slush (Large amounts of snow on a wet surface)

7 - Light or patchy snow (Snow on surface where road is also visible but not wet)

8 - Heavy snow (Snow on surface not visible road)

9 - Obstructed. Cannot decide class


## Example annotations . Click to enlarge
### 0 - Dry (No visible water or patches on the surface. Surrounding environment also looks dryish)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/05/02/250/250_20230502T1200Z.jpg&cc=0&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/05/02/250/250_20230502T1200Z.jpg&cc=0&obs2=-1)
[<img  height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/05/02/27/27_20230502T0600Z.jpg&cc=0&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/05/02/27/27_20230502T0600Z.jpg&cc=0&obs2=-1)
[<img  height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/26/27/27_20230426T1200Z.jpg&cc=0&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/26/27/27_20230426T1200Z.jpg&cc=0&obs2=-1)

<br/>


### 1 - Patchy water (Patches of water or not dry looking, with moistness in the surface)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/26/250/250_20230426T1800Z.jpg&cc=1&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/26/250/250_20230426T1800Z.jpg&cc=1&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/30/149/149_20230430T0600Z.jpg&cc=1&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/30/149/149_20230430T0600Z.jpg&cc=1&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/09/149/149_20230409T1200Z.jpg&cc=1&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/09/149/149_20230409T1200Z.jpg&cc=1&obs2=-1)

<br/>

### 2 - Light water (Thinner film of water, light reflection, often wet or snowy surroundings)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/24/27/27_20230424T1200Z.jpg&cc=2&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/24/27/27_20230424T1200Z.jpg&cc=2&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/05/01/27/27_20230501T1200Z.jpg&cc=2&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/05/01/27/27_20230501T1200Z.jpg&cc=2&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/03/02/536/536_20230302T0600Z.jpg&cc=2&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/03/02/536/536_20230302T0600Z.jpg&cc=2&obs2=-1)

<br/>

### 3 - Heavy water (Clearly very wet)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/25/27/27_20230425T0000Z.jpg&cc=3&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/25/27/27_20230425T0000Z.jpg&cc=3&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/21/149/149_20230421T1800Z.jpg&cc=3&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/21/149/149_20230421T1800Z.jpg&cc=3&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/24/536/536_20230424T1800Z.jpg&cc=3&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/24/536/536_20230424T1800Z.jpg&cc=3&obs2=-1)

<br/>

### 4 - Ice
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/02/04/23/23_20230204T1200Z.jpg&cc=4&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/02/04/23/23_20230204T1200Z.jpg&cc=4&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/02/403/403_20230402T1200Z.jpg&cc=4&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/02/403/403_20230402T1200Z.jpg&cc=4&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/03/13/248/248_20230313T0000Z.jpg&cc=4&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/03/13/248/248_20230313T0000Z.jpg&cc=4&obs2=-1)

<br/>

### 5 - Light slush (Small amounts of snow on a wet surface)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/24/149/149_20230424T1800Z.jpg&cc=5&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/24/149/149_20230424T1800Z.jpg&cc=5&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/04/149/149_20230404T1200Z.jpg&cc=5&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/04/149/149_20230404T1200Z.jpg&cc=5&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/02/04/65/65_20230204T1200Z.jpg&cc=5&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/02/04/65/65_20230204T1200Z.jpg&cc=5&obs2=-1)

<br/>

### 6 - Heavy slush (Large amounts of snow on a wet surface)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/22/149/149_20230422T1200Z.jpg&cc=6&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/22/149/149_20230422T1200Z.jpg&cc=6&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/02/07/49/49_20230207T1800Z.jpg&cc=6&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/02/07/49/49_20230207T1800Z.jpg&cc=6&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/02/06/244/244_20230206T1200Z.jpg&cc=6&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/02/06/244/244_20230206T1200Z.jpg&cc=6&obs2=-1)

<br/>

### 7 - Light or patchy snow (Snow on surface where road is also visible but not wet)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/30/250/250_20230430T0600Z.jpg&cc=7&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/30/250/250_20230430T0600Z.jpg&cc=7&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/04/149/149_20230404T1800Z.jpg&cc=7&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/04/149/149_20230404T1800Z.jpg&cc=7&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/03/20/149/149_20230320T1200Z.jpg&cc=7&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/03/20/149/149_20230320T1200Z.jpg&cc=7&obs2=-1)

<br/>

### 8 - Heavy snow (Snow on surface not visible road)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/23/149/149_20230423T1200Z.jpg&cc=8&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/23/149/149_20230423T1200Z.jpg&cc=8&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/03/149/149_20230403T0000Z.jpg&cc=8&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/04/03/149/149_20230403T0000Z.jpg&cc=8&obs2=-1)
[<img height="224" width="224" src="https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/03/31/149/149_20230331T0000Z.jpg&cc=8&obs2=-1">](https://modellprod.k8s.met.no/roadlabels/labeledimage?q=roadcams/2023/03/31/149/149_20230331T0000Z.jpg&cc=8&obs2=-1)

<br/>

See more examples in the annotation app 


## Bugs And Restrictions
- mobile screen size not supported. (But seems possible to annotate still)
- We are using some new (for me at least) stacks at met. Kubernetes and objectstore . One of them is slow at times . Investigation continues.
