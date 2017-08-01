# Stream Conditions

Stream Conditions allow you to define conditions if a Stream should be included in a Encoding or not.
More information about why Stream Conditions are useful can be found in [our Blogpost](https://bitmovin.com/stream-conditions-video-encoding-workflows/).

## Usage

StreamConditions are part of the `models.Stream` Object that gets sent to `POST /v1/encoding/encodings/<encodingId>/streams` and get created through the following Methods:

* `models.NewAttributeCondition` 
* `models.NewAndConjunction`
* `models.NewOrDisjunction`

## Attribute Conditions

As the name implies `models.AttributeConditions` represent a Condition check on an attribute, with an Operator and a Value that gets checked. (Think of it as a if statement)

Declaration:

```
NewAttributeCondition(attribute bitmovintypes.ConditionAttribute, operator, value string) StreamCondition
```

The following condition translated to:

`models.NewAttributeCondition(bitmovintypes.ConditionAttributeHeight, ">=", "1080")
`

> _if Input.Height >= 1080_

## Conjunction and Disjunction

You can also define logical Conjunctions and Disjunctions (AND and OR) by chaining the Conditions together:

```
models.NewAndConjunction(
  models.NewAttributeCondition(bitmovintypes.ConditionAttributeHeight, "==", "1080"),
  models.NewAttributeCondition(bitmovintypes.ConditionAttributeWidth, "==", "1920"),
)
```

This example will only encode the Stream if it's a real 16:9 FullHD Input File and translates to 

> _If Input.Height == 1080 && Input.Width == 1920_

The same goes for the `models.NewOrDisjunction` just that it represents a logical **OR** of all parameters.

## Full Example

A full usage example with a Stream would look like this:

```
inputStream := models.InputStream{
  InputID:        inputID, // ID of the Input Resource (bucket)
  InputPath:      stringToPtr("/path/to/file.mp4"),
  SelectionMode:  bitmovinTypes.SelectionModeAuto,
}
videoStream1080p := &models.Stream{
  CodecConfigurationID: 1080ConfigurationID, //Existing CodecConfiguration resource ID
  InputStreams:         inputStream,
  Conditions:           models.NewAttributeCondition(bitmovintypes.ConditionAttributeHeight, ">=", "1080"),
}
encodingService.AddStream(*encodingID, videoStream1080p) // This stream will only be encoded if the Input width is >= 1080
```