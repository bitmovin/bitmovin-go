package bitmovintypes

type FragmentedMP4MuxingManifestType string

const (
	FragmentedMP4MuxingManifestTypeSmooth                      FragmentedMP4MuxingManifestType = "SMOOTH"
	FragmentedMP4MuxingManifestTypeDASHOnDemand                FragmentedMP4MuxingManifestType = "DASH_ON_DEMAND"
	FragmentedMP4MuxingManifestTypeHlsByteRanges               FragmentedMP4MuxingManifestType = "HLS_BYTE_RANGES"
	FragmentedMP4MuxingManifestTypeHlsByteRangesIFramePlaylist FragmentedMP4MuxingManifestType = "HLS_BYTE_RANGES_AND_IFRAME_PLAYLIST"
)
