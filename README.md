<p align="center">
  <img src="/assets/logo.svg" width="250px">
</p>

<p align="center">An iterative algorithm to generate high quality triangulated and polygonal art from images.</p>
<p align="center">
<a><img src="https://github.com/RH12503/Triangula/actions/workflows/test.yml/badge.svg" alt="Test status"></a>
<a href="https://pkg.go.dev/github.com/RH12503/Triangula"><img src="https://pkg.go.dev/badge/github.com/RH12503/Triangula.svg" alt="Go Reference"></a>
<a href="https://goreportcard.com/report/github.com/RH12503/Triangula"><img src="https://goreportcard.com/badge/github.com/RH12503/Triangula" alt="Go Report Card"></a>
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
<a href="https://twitter.com/intent/tweet?text=An%20iterative%20algorithm%20to%20triangulate%20images.&url=https://github.com/RH12503/triangula&hashtags=golang,geneticalgorithm,generativeart"><img src="https://img.shields.io/twitter/url/http/shields.io.svg?style=social" alt="Tweet"></a>
</p>

Triangula uses a modified genetic algorithm to triangulate or polygonate images. It works best with images smaller than 3000px and with fewer than 3000 points, typically producing an optimal result within a couple of minutes. For a full explanation of the algorithm, see [this page in the wiki](https://github.com/RH12503/Triangula/wiki/Explanation-of-the-algorithm). 

You can try the algorithm out in your browser [here](https://rh12503.github.io/triangula/), but the desktop app will typically be 20-50x faster. 

## Install

### GUI
Install the [GUI](https://github.com/RH12503/Triangula-GUI) from the [releases page](https://github.com/RH12503/Triangula/releases). 
The GUI uses [Wails](https://wails.app/) for its frontend. 
<p float="left" align="middle">
<img src="/assets/gui1.png" width="49%">
<img src="/assets/gui2.png" width="49%">
</p>
If the app isn't running on Linux, go to the Permissions tab in the executable's properties and tick `Allow executing file as program`. 

### CLI
Install the [CLI](https://github.com/RH12503/Triangula-CLI) by running: 
```
go get -u github.com/RH12503/Triangula-CLI/triangula
```

Your `PATH` variable also needs to include your `go/bin` directory, which is `~/go/bin` on macOS, `$GOPATH/bin` on Linux, and `c:\Go\bin` on Windows. 

Then run it using the command: 
```
triangula run -img <path to image> -out <path to output JSON>
```

and when you're happy with its fitness, render a SVG:
```
triangula render -in <path to outputted JSON> -img <path to image> -out <path to output SVG> 
```
For more detailed instructions, including rendering PNGs with effects see [this page](https://github.com/RH12503/Triangula-CLI/blob/main/README.md). 

## Options
For almost all cases, only changing the number of points and leaving all other options with their default values will generate an optimal result. 

| Name  | Flag | Default |  Usage |
| ------------- | ---- | ------------- | -- |
|  Points |  `--points, -p`  | 300 | The number of points to use in the triangulation   |
| Mutations  |  `--mutations, --mut, -m`  |  2 | The number of mutations to make |
| Variation | `--variation, -v` | 0.3 | The variation each mutation causes |
| Population | `--population, --pop, --size` | 400 | The population size in the algorithm |
| Cutoff | `--cutoff, --cut` | 5 | The cutoff value of the algorithm |
| Cache | `--cache, -c` | 22 | The cache size as a power of 2 |
| Block | `--block, -b` | 5 | The size of the blocks used when rendering |
| Threads | `--threads, -t` | 0 | The number of threads to use or 0 to use all cores | 
| Repetitions | `--reps, -r`| 500 | The number of generations before saving to the output file (CLI only) | 

## Examples of output
### Triangulated
<img src="/assets/output/grad.png" height="400"/>
<img src="/assets/output/plane.png" height="400"/> 
<img src="/assets/output/sf.png" height="400"/>
<img src="/assets/output/elon.png" height="400"/>
<img src="/assets/output/astro.png" height="400"/>

### Polygonal
<img src="/assets/output/dog.png" height="400"/>
<img src="/assets/output/obama.png" height="400"/> 
<img src="/assets/output/science.png" height="400"/>
<img src="/assets/output/queen.png" height="400"/>

### Comparison to [esimov/triangle](https://github.com/esimov/triangle)
esimov/triangle seems to be a similar project to Triangula that is also written in Go. However, the two appear to generate very different styles. One big advantage of triangle is that it generates an image almost instantaneously, while Triangula needs to run many iterations. 

esimov/triangle results were taken from their [Github repo](https://github.com/esimov/triangle), and Triangula's results were generated over 1-2 minutes. 
| esimov/triangle | Triangula |
| :---: | :---: |
| <img src="https://github.com/esimov/triangle/blob/master/output/sample_11.png" height="250"/> | <img src="/assets/output/result.png" height="250"/> |
| <img src="https://github.com/esimov/triangle/blob/master/output/sample_3.png" height="250"/> | <img src="/assets/output/result2.png" height="250"/>  |

#### Difference from [fogleman/primitive](https://github.com/fogleman/primitive) and [gheshu/image_decompiler](https://github.com/gheshu/image_decompiler)
A lot of people have commented about Triangula's similarities to these other algorithms. While all these algorithms are iterative algorithms, the main difference is that in the other algorithms triangles can overlap while Triangula generates a triangulation. 

## API 
Simple example: 
```Go
import imageData "github.com/RH12503/Triangula/image"

func main() {
    // Open and decode a PNG/JPEG
    file, err := os.Open("image.png")

    if err != nil {
          log.Fatal(err)
    }

    image, _, err := image.Decode(file)

    file.Close()

    if err != nil {
          log.Fatal(err)
    }

    img := imageData.ToData(image)


    pointFactory := func() normgeom.NormPointGroup {
          return (generator.RandomGenerator{}).Generate(200) // 200 points
    }

    evaluatorFactory := func(n int) evaluator.Evaluator {
          // 22 for the cache size and 5 for the block size
          // use PolygonsImageFunctions for polygons 
		  return evaluator.NewParallel(fitness.TrianglesImageFunctions(imgData, 5, n), 22)
    }

    var mutator mutation.Method

    // 1% mutation rate and 30% variation
    mutator = mutation.NewGaussianMethod(0.01, 0.3)

    // 400 population size and 5 cutoff
    algo := algorithm.NewModifiedGenetic(pointFactory, 400, 5, evaluatorFactory, mutator)

    // Run the algorithm
    for {
          algo.Step()
          fmt.Println(algo.Stats().BestFitness)
    }
}

```
## Contribute
Any contributions are welcome. Currently help is needed with:
* Support for exporting effects to SVGs. 
* Supporting more image types for the CLI and GUI. (eg. .tiff, .webp, .heic)
* Allowing drag and drop of images from the web for the GUI. 
* More effects. 
* Any optimizations. 

Thank you to these contributors for helping to improve Triangula: 
- [@bstncartwright](https://github.com/bstncartwright)
- [@2BoysAndHats](https://github.com/2BoysAndHats)
