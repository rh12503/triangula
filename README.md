<p align="center">
  <img src="/assets/logo.svg" width="250px">
</p>

<p align="center">An iterative algorithm to generate high quality triangulated images.</p>
<p align="center">
<a><img src="https://github.com/RH12503/Triangula/actions/workflows/test.yml/badge.svg" alt="Test status"></a>
<a href="https://pkg.go.dev/github.com/RH12503/Triangula"><img src="https://pkg.go.dev/badge/github.com/RH12503/Triangula.svg" alt="Go Reference"></a>
<a href="https://goreportcard.com/report/github.com/RH12503/Triangula"><img src="https://goreportcard.com/badge/github.com/RH12503/Triangula" alt="Go Report Card"></a>
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
</p>

## Introduction
Triangula uses a modified genetic algorithm to triangulate images. It works best with images smaller than 2000px and with less than 3000 points, typically producing an optimal result within a couple of minutes. To better understand the algorithm, considering reading [this page on the wiki](https://github.com/RH12503/Triangula/wiki/Explanation-of-the-algorithm). 

## Install

### GUI
Install the [GUI](https://github.com/RH12503/Triangula-GUI) from the releases page. 

<img src="/assets/triangula.gif" width="500">

### CLI
Install the [CLI](https://github.com/RH12503/Triangula-CLI) using: 
```
go get github.com/RH12503/Triangula-CLI/triangula
```

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

## Example output
<img src="/assets/output/grad.png" height="400"/>
<img src="/assets/output/plane.png" height="400"/> 
<img src="/assets/output/sf.png" height="400"/>
<img src="/assets/output/thompson.png" height="400"/>
<img src="/assets/output/astro.png" height="400"/>

### Comparison to [esimov/triangle](https://github.com/esimov/triangle)
esimov/triangle seems to be a similar project to Triangula that is also written in Go. However, the two appear to generate very different styles. One big advantage of triangle is that it generates an image almost instantaneously, while Triangula needs to run many iterations. 
| esimov/triangle | Triangula |
| :---: | :---: |
| <img src="https://github.com/esimov/triangle/blob/master/output/sample_11.png" height="250"/> | <img src="/assets/output/result.png" height="250"/> |
| <img src="https://github.com/esimov/triangle/blob/master/output/sample_3.png" height="250"/> | <img src="/assets/output/result2.png" height="250"/>  |

## API 
A simple example to use the API would be: 
```Go
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
            return evaluator.NewParallel(img, 22, 5, n)
      }

      var mutator mutation.Method

      // 1% mutation rate and 30% variation
      mutator = mutation.NewGaussianMethod(0.01, 0.3)

      // 400 population size and 5 cutoff
      algo := algorithm.NewSimple(pointFactory, 400, 5, evaluatorFactory, mutator)

      // Run the algorithm
      for {
            algo.Step()
            fmt.Println(algo.Stats().BestFitness)
      }
}

```
