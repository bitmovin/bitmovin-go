# Creating Thumbnails

Thumbnails can be extracted from a Video on a per-stream basis. So you can extract a 1080p Thumbnail from your 1080p Stream while generating a smaller thumbnail for the 720p representation.

A full example with Thumbnails is available in the Repository at [examples/simple_encoding_thumbnails/main.go](https://github.com/streamco/bitmovin-go/blob/master/examples/simple_encoding_thumbnails/main.go).

Thumbnails have 3 mandatory properties (height, positions and output) and are accessible through the `models.NewThumbnail()` constructor (or directly through the `models.Thumbnail` struct).

This example assumes you already created an encoding and a 1080p stream:

Create an output object:

```go
thumbOutput := &models.Output{
  OutputID:   outputResponse.Data.Result.ID, // You can reuse the output for the Stream here
  OutputPath: stringToPtr("thumbs/"), // The path where you want your thumbnails
}
```


Next create the Thumbnail

```go
// This will generate one thumbnail for second 3, second 5 and second 30 in the video
thumb1080 := models.NewThumbnail(400, []float64{3, 5, 30}, []models.Output{*thumbOutput})

if _, err := encodingS.AddThumbnail(encodingID, 1080VideoStreamID, thumb1080); err != nil {
  log.Fatalf("Error creating 1080p Thumbnail resource")
}

```

As you can see from the definition of `NewThumbnail` you can specify multiple outputs to write the Thumbnail to. To do so just add additional `model.Output` models to `[]models.Output{}`.

## Changing the Position Unit

Multiple Thumbnails can be generated during a single stream encoding. That's why the `position` parameter takes an array of `float64`. By default the specified positions are seconds. You can also change this to take screen grabs at a percentage of the video duration:

```go
thumb := models.NewThumbnail(400, []float64{10}, []models.Output{*thumbOutput}).Builder().
  PositionUnit(bitmovintypes.bitmovintypes.PositionPercents).Build()
```


## Changing the Filename Pattern

By default the pattern for Thumbnails will be `thumbnail-%number%.png`. If you want a different naming scheme for your Thumbnails you can specify it through the `Pattern` property.

```go
thumb := models.NewThumbnail(400, []float64{30}, []models.Output{*thumbOutput}).Builder().
  Pattern("video-thumb-%number%.png").Build()
```

## Fluent Builder

The Thumbnail model uses the new Fluent Builder API that avoids the `stringToPtr()` to set nullable values. If you dislike it you can still use the regular struct based syntax:

```go
thumb2 := &models.Thumbnail{
  Height:       400,
  Positions:    []float64{30},
  PositionUnit: bitmovintypes.PositionPercents,
  Outputs:      []models.Output{*thumbOutput},
  Pattern:      stringToPtr("video-thumb-%number%.png"),
}
```

The new Fluent Builder enables a bit cleaner syntax:

```go
thumb := models.NewThumbnail(400, []float64{30}, []models.Output{*thumbOutput}).Builder().
  PositionUnit(bitmovintypes.PositionPercents).
  Pattern("video-thumb-%number%.png").
  Build()
```