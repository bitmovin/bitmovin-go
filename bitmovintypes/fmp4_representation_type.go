package bitmovintypes

type FMP4RepresentationType string

const (
	FMP4RepresentationTypeTemplate              FMP4RepresentationType = "TEMPLATE"
	FMP4RepresentationTypeList                  FMP4RepresentationType = "LIST"
	FMP4RepresentationTypeTimeline              FMP4RepresentationType = "TIMELINE"
	FMP4RepresentationTypeTemplateAdaptationSet FMP4RepresentationType = "TEMPLATE_ADAPTATION_SET"
)
