package slack

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDataVisualizationBlock(t *testing.T) {
	chart := NewDataVisualizationPieChart(
		NewDataVisualizationSegment("Kit Kat", 45),
		NewDataVisualizationSegment("Twix", 28),
	)
	block := NewDataVisualizationBlock("Candy Bars", chart,
		DataVisualizationBlockOptionBlockID("dv-1"),
	)

	assert.Equal(t, MBTDataVisualization, block.BlockType())
	assert.Equal(t, "data_visualization", string(block.Type))
	assert.Equal(t, "dv-1", block.ID())
	assert.Equal(t, "Candy Bars", block.Title)
	assert.Equal(t, DataVisualizationChartPie, block.Chart.DataVisualizationChartType())
}

func TestNewDataVisualizationBlockWithNilOption(t *testing.T) {
	assert.NotPanics(t, func() {
		NewDataVisualizationBlock("title", NewDataVisualizationPieChart(), nil)
	}, "should not panic when nil option passed")
}

func TestDataVisualizationBlockPieJSONRoundTrip(t *testing.T) {
	payload := `{
		"type": "data_visualization",
		"block_id": "dv-pie",
		"title": "My Favorite Candy Bars",
		"chart": {
			"type": "pie",
			"segments": [
				{"label": "Kit Kat", "value": 45},
				{"label": "Twix", "value": 28},
				{"label": "Crunch", "value": 18},
				{"label": "Milky Way", "value": 9}
			]
		}
	}`

	var block DataVisualizationBlock
	require.NoError(t, json.Unmarshal([]byte(payload), &block))

	assert.Equal(t, MBTDataVisualization, block.BlockType())
	assert.Equal(t, "dv-pie", block.ID())
	assert.Equal(t, "My Favorite Candy Bars", block.Title)
	pie, ok := block.Chart.(*DataVisualizationPieChart)
	require.True(t, ok)
	require.Len(t, pie.Segments, 4)
	assert.Equal(t, "Kit Kat", pie.Segments[0].Label)
	assert.Equal(t, float64(45), pie.Segments[0].Value)

	marshalled, err := json.Marshal(block)
	require.NoError(t, err)

	var expected, actual map[string]any
	require.NoError(t, json.Unmarshal([]byte(payload), &expected))
	require.NoError(t, json.Unmarshal(marshalled, &actual))
	assert.Equal(t, expected, actual)
}

func TestDataVisualizationBlockSeriesChartsJSONRoundTrip(t *testing.T) {
	payload := `[
		{
			"type": "data_visualization",
			"block_id": "dv-bar",
			"title": "Pie Tastiness",
			"chart": {
				"type": "bar",
				"series": [
					{
						"name": "Pies",
						"data": [
							{"label": "Pumpkin", "value": 70},
							{"label": "Blueberry", "value": 90}
						]
					}
				],
				"axis_config": {
					"categories": ["Pumpkin", "Blueberry"],
					"x_label": "Pies",
					"y_label": "Percentage of Tastiness"
				}
			}
		},
		{
			"type": "data_visualization",
			"block_id": "dv-area",
			"title": "Daily Active Users",
			"chart": {
				"type": "area",
				"series": [
					{
						"name": "Free Tier",
						"data": [
							{"label": "Mon", "value": 12000},
							{"label": "Tue", "value": 13500}
						]
					},
					{
						"name": "Paid Tier",
						"data": [
							{"label": "Mon", "value": 4500},
							{"label": "Tue", "value": 4800}
						]
					}
				],
				"axis_config": {
					"categories": ["Mon", "Tue"],
					"x_label": "Day",
					"y_label": "Users"
				}
			}
		},
		{
			"type": "data_visualization",
			"block_id": "dv-line",
			"title": "Weekly Paper Sales",
			"chart": {
				"type": "line",
				"series": [
					{
						"name": "Website",
						"data": [
							{"label": "Week 1", "value": 32000},
							{"label": "Week 2", "value": 35000}
						]
					}
				],
				"axis_config": {
					"categories": ["Week 1", "Week 2"],
					"x_label": "Week",
					"y_label": "Paper Sales (USD)"
				}
			}
		}
	]`

	var blocks Blocks
	require.NoError(t, json.Unmarshal([]byte(payload), &blocks))
	require.Len(t, blocks.BlockSet, 3)

	bar, ok := blocks.BlockSet[0].(*DataVisualizationBlock)
	require.True(t, ok, "expected *DataVisualizationBlock, got %T", blocks.BlockSet[0])
	assert.Equal(t, MBTDataVisualization, bar.BlockType())
	assert.Equal(t, "dv-bar", bar.ID())
	barChart, ok := bar.Chart.(*DataVisualizationBarChart)
	require.True(t, ok)
	require.Len(t, barChart.Series, 1)
	assert.Equal(t, "Pies", barChart.Series[0].Name)
	assert.Equal(t, []string{"Pumpkin", "Blueberry"}, barChart.AxisConfig.Categories)

	area := blocks.BlockSet[1].(*DataVisualizationBlock)
	areaChart, ok := area.Chart.(*DataVisualizationAreaChart)
	require.True(t, ok)
	require.Len(t, areaChart.Series, 2)
	assert.Equal(t, DataVisualizationChartArea, areaChart.DataVisualizationChartType())

	line := blocks.BlockSet[2].(*DataVisualizationBlock)
	lineChart, ok := line.Chart.(*DataVisualizationLineChart)
	require.True(t, ok)
	assert.Equal(t, "Week", lineChart.AxisConfig.XLabel)
	assert.Equal(t, "Paper Sales (USD)", lineChart.AxisConfig.YLabel)

	marshalled, err := json.Marshal(blocks)
	require.NoError(t, err)

	var expected, actual []map[string]any
	require.NoError(t, json.Unmarshal([]byte(payload), &expected))
	require.NoError(t, json.Unmarshal(marshalled, &actual))
	assert.Equal(t, expected, actual)
}

func TestDataVisualizationBlockConstructorsMarshal(t *testing.T) {
	axis := NewDataVisualizationAxisConfig("Week 1", "Week 2").
		WithXLabel("Week").
		WithYLabel("Paper Sales")
	series := NewDataVisualizationDataSeries("Website",
		NewDataVisualizationDataPoint("Week 1", 32000),
		NewDataVisualizationDataPoint("Week 2", 35000),
	)
	block := NewDataVisualizationBlock(
		"Weekly Paper Sales",
		NewDataVisualizationLineChart(axis, series),
	)

	marshalled, err := json.Marshal(block)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"type": "data_visualization",
		"title": "Weekly Paper Sales",
		"chart": {
			"type": "line",
			"series": [
				{
					"name": "Website",
					"data": [
						{"label": "Week 1", "value": 32000},
						{"label": "Week 2", "value": 35000}
					]
				}
			],
			"axis_config": {
				"categories": ["Week 1", "Week 2"],
				"x_label": "Week",
				"y_label": "Paper Sales"
			}
		}
	}`, string(marshalled))
}

func TestDataVisualizationBlockUnknownChartType(t *testing.T) {
	payload := `{
		"type": "data_visualization",
		"title": "Mystery Chart",
		"chart": {
			"type": "scatter",
			"series": []
		}
	}`

	var block DataVisualizationBlock
	err := json.Unmarshal([]byte(payload), &block)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "scatter")
}

func TestDataVisualizationBlockValidate(t *testing.T) {
	validPie := NewDataVisualizationBlock("Candy Bars",
		NewDataVisualizationPieChart(NewDataVisualizationSegment("Kit Kat", 45)),
	)
	require.NoError(t, validPie.Validate())

	validLine := NewDataVisualizationBlock("Weekly Paper Sales",
		NewDataVisualizationLineChart(
			NewDataVisualizationAxisConfig("Week 1", "Week 2").
				WithXLabel("Week").
				WithYLabel("Paper Sales"),
			NewDataVisualizationDataSeries("Website",
				NewDataVisualizationDataPoint("Week 1", 32000),
				NewDataVisualizationDataPoint("Week 2", -10),
			),
		),
	)
	require.NoError(t, validLine.Validate(), "negative data point values are allowed")
}

func TestDataVisualizationBlockValidateFieldLimits(t *testing.T) {
	tests := []struct {
		name      string
		block     *DataVisualizationBlock
		errorText string
	}{
		{
			name: "missing title",
			block: NewDataVisualizationBlock("",
				NewDataVisualizationPieChart(NewDataVisualizationSegment("Kit Kat", 45)),
			),
			errorText: "title must have a minimum length of 1",
		},
		{
			name: "title too long",
			block: NewDataVisualizationBlock(strings.Repeat("a", 51),
				NewDataVisualizationPieChart(NewDataVisualizationSegment("Kit Kat", 45)),
			),
			errorText: "title cannot be longer than 50 characters",
		},
		{
			name: "missing chart",
			block: &DataVisualizationBlock{
				Type:  MBTDataVisualization,
				Title: "No Chart",
			},
			errorText: "chart is required",
		},
		{
			name: "typed nil chart",
			block: &DataVisualizationBlock{
				Type:  MBTDataVisualization,
				Title: "Typed Nil Chart",
				Chart: (*DataVisualizationLineChart)(nil),
			},
			errorText: "chart is required",
		},
		{
			name: "segment label too long",
			block: NewDataVisualizationBlock("Candy Bars",
				NewDataVisualizationPieChart(NewDataVisualizationSegment(strings.Repeat("a", 21), 45)),
			),
			errorText: "segment 0 label cannot be longer than 20 characters",
		},
		{
			name: "segment value is zero",
			block: NewDataVisualizationBlock("Candy Bars",
				NewDataVisualizationPieChart(NewDataVisualizationSegment("Kit Kat", 0)),
			),
			errorText: "segment 0 value must be greater than 0",
		},
		{
			name: "series name too long",
			block: NewDataVisualizationBlock("Weekly Paper Sales",
				NewDataVisualizationLineChart(
					NewDataVisualizationAxisConfig("Week 1"),
					NewDataVisualizationDataSeries(strings.Repeat("a", 21),
						NewDataVisualizationDataPoint("Week 1", 32000),
					),
				),
			),
			errorText: "series 0 name cannot be longer than 20 characters",
		},
		{
			name: "category too long",
			block: NewDataVisualizationBlock("Weekly Paper Sales",
				NewDataVisualizationLineChart(
					NewDataVisualizationAxisConfig(strings.Repeat("a", 21)),
					NewDataVisualizationDataSeries("Website",
						NewDataVisualizationDataPoint(strings.Repeat("a", 21), 32000),
					),
				),
			),
			errorText: "axis_config.categories[0] cannot be longer than 20 characters",
		},
		{
			name: "axis label too long",
			block: NewDataVisualizationBlock("Weekly Paper Sales",
				NewDataVisualizationLineChart(
					NewDataVisualizationAxisConfig("Week 1").WithXLabel(strings.Repeat("a", 51)),
					NewDataVisualizationDataSeries("Website",
						NewDataVisualizationDataPoint("Week 1", 32000),
					),
				),
			),
			errorText: "axis_config.x_label cannot be longer than 50 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.block.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorText)
		})
	}
}

func TestDataVisualizationBlockValidatePieConstraints(t *testing.T) {
	noSegments := NewDataVisualizationBlock("Candy Bars", NewDataVisualizationPieChart())
	err := noSegments.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pie chart must have at least 1 segment")

	tooManySegments := NewDataVisualizationBlock("Candy Bars",
		NewDataVisualizationPieChart(
			NewDataVisualizationSegment("A", 1),
			NewDataVisualizationSegment("B", 1),
			NewDataVisualizationSegment("C", 1),
			NewDataVisualizationSegment("D", 1),
			NewDataVisualizationSegment("E", 1),
			NewDataVisualizationSegment("F", 1),
			NewDataVisualizationSegment("G", 1),
		),
	)
	err = tooManySegments.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pie chart cannot have more than 6 segments")
}

func TestDataVisualizationBlockValidateSeriesRuntimeRules(t *testing.T) {
	tests := []struct {
		name      string
		block     *DataVisualizationBlock
		errorText string
	}{
		{
			name: "missing axis config categories",
			block: NewDataVisualizationBlock("Weekly Paper Sales",
				NewDataVisualizationLineChart(DataVisualizationAxisConfig{},
					NewDataVisualizationDataSeries("Website",
						NewDataVisualizationDataPoint("Week 1", 32000),
					),
				),
			),
			errorText: "axis_config.categories must have at least 1 category",
		},
		{
			name: "data point label outside categories",
			block: NewDataVisualizationBlock("Weekly Paper Sales",
				NewDataVisualizationLineChart(
					NewDataVisualizationAxisConfig("Week 1"),
					NewDataVisualizationDataSeries("Website",
						NewDataVisualizationDataPoint("Week 2", 32000),
					),
				),
			),
			errorText: `label "Week 2" must match axis_config.categories`,
		},
		{
			name: "omitted category point",
			block: NewDataVisualizationBlock("Weekly Paper Sales",
				NewDataVisualizationLineChart(
					NewDataVisualizationAxisConfig("Week 1", "Week 2"),
					NewDataVisualizationDataSeries("Website",
						NewDataVisualizationDataPoint("Week 1", 32000),
					),
				),
			),
			errorText: "series 0 data must contain exactly one point for every category",
		},
		{
			name: "duplicate data label",
			block: NewDataVisualizationBlock("Weekly Paper Sales",
				NewDataVisualizationLineChart(
					NewDataVisualizationAxisConfig("Week 1", "Week 2"),
					NewDataVisualizationDataSeries("Website",
						NewDataVisualizationDataPoint("Week 1", 32000),
						NewDataVisualizationDataPoint("Week 1", 35000),
					),
				),
			),
			errorText: `series 0 data must not contain duplicate label "Week 1"`,
		},
		{
			name: "duplicate series name",
			block: NewDataVisualizationBlock("Weekly Paper Sales",
				NewDataVisualizationLineChart(
					NewDataVisualizationAxisConfig("Week 1"),
					NewDataVisualizationDataSeries("Website",
						NewDataVisualizationDataPoint("Week 1", 32000),
					),
					NewDataVisualizationDataSeries("Website",
						NewDataVisualizationDataPoint("Week 1", 35000),
					),
				),
			),
			errorText: `series names must be unique: "Website"`,
		},
		{
			name: "too many series",
			block: NewDataVisualizationBlock("Weekly Paper Sales",
				NewDataVisualizationLineChart(
					NewDataVisualizationAxisConfig("Week 1"),
					NewDataVisualizationDataSeries("A", NewDataVisualizationDataPoint("Week 1", 1)),
					NewDataVisualizationDataSeries("B", NewDataVisualizationDataPoint("Week 1", 1)),
					NewDataVisualizationDataSeries("C", NewDataVisualizationDataPoint("Week 1", 1)),
					NewDataVisualizationDataSeries("D", NewDataVisualizationDataPoint("Week 1", 1)),
					NewDataVisualizationDataSeries("E", NewDataVisualizationDataPoint("Week 1", 1)),
					NewDataVisualizationDataSeries("F", NewDataVisualizationDataPoint("Week 1", 1)),
					NewDataVisualizationDataSeries("G", NewDataVisualizationDataPoint("Week 1", 1)),
				),
			),
			errorText: "line chart cannot have more than 6 series",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.block.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorText)
		})
	}
}
