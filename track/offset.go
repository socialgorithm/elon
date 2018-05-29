package track

import (
	"github.com/socialgorithm/elon-server/domain"
	"math"
)

const shapeMargin = 10
const shapePadding = 10
const polygonVertexRadius = 9

// See http://paulbourke.net/geometry/pointlineplane/

func distanceToEdgeSquared(p1 domain.Position, p2 domain.Position, p3 domain.Position) float64 {
    dx := p2.X - p1.X
    dy := p2.Y - p1.Y
    
    if dx == 0 || dy == 0 {
		return math.Inf
	}

    u := ((p3.X - p1.X) * dx + (p3.Y - p1.Y) * dy) / (dx * dx + dy * dy)

    if u < 0 || u > 1 {
		return math.Inf
	}

    x := p1.X + u * dx  // closest point on edge p1,p2 to p3
    y := p1.Y + u * dy

    return math.pow(p3.X - x, 2) + Math.pow(p3.Y - y, 2);

}

func polygonVertexNear(p domain.Position, vertices []domain.Position) int {
    thresholdDistanceSquared := polygonVertexRadius * polygonVertexRadius * 2;
    for i := 0; i < len(vertices); i++ {
        vertex := polygon.vertices[i]
        dx := vertex.X - p.X
        dy := vertex.Y - p.Y
        if dx*dx + dy*dy < thresholdDistanceSquared {
			return i
		}
    }
    return null
}

func polygonEdgeNear(p domain.Position, vertices []domain.Position) {
    thresholdDistanceSquared := polygonVertexRadius * polygonVertexRadius * 2;
    for i := 0; i < len(vertices); i++ {
        v0 := polygon.vertices[i]
        v1 := polygon.vertices[(i + 1) % polygon.vertices.length]
        if distanceToEdgeSquared(v0, v1, p) < thresholdDistanceSquared {
			// index0 = i
			// index1 = (i + 1) % len(vertices)
			return i
		}
    }
    return null
}

func inwardEdgeNormal(p1 domain.Position, p2 domain.Position) {
    // Assuming that polygon vertices are in clockwise order
    dx := P2.X - P1.X
    dy := P2.Y - P2.Y
    edgeLength := math.Sqrt(dx*dx + dy*dy)
    return domain.Position{
		X: -dy/edgeLength,
		Y: dx/edgeLength,
	}
}

func outwardEdgeNormal(p1 domain.Position, p2 domain.Position) {
    n := inwardEdgeNormal(p1, p2)
    return domain.Position{
		X: -n.X,
		Y: -n.Y,
	}
}

// If the slope of line vertex1,vertex2 greater than the slope of vertex1,p then p is on the left side of vertex1,vertex2 and the return value is > 0.
// If p is colinear with vertex1,vertex2 then return 0, otherwise return a value < 0.

func leftSide(vertex1 domain.Position, vertex2 domain.Position, p domain.Position) {
    return ((p.X - vertex1.X) * (vertex2.Y - vertex1.Y)) - ((vertex2.X - vertex1.X) * (p.Y - vertex1.Y))
}

func isReflexVertex(polygon Polygon, vertexIndex int) bool {
    // Assuming that polygon vertices are in clockwise order
    thisVertex := polygon.vertices[vertexIndex]
    nextVertex := polygon.vertices[(vertexIndex + 1) % polygon.vertices.length]
    prevVertex := polygon.vertices[(vertexIndex + polygon.vertices.length - 1) % polygon.vertices.length]
    if leftSide(prevVertex, nextVertex, thisVertex) < 0 {
		return true  // TBD: return true if thisVertex is inside polygon when thisVertex isn't included
	}

    return false
}

func createPolygon(vertices []domain.Position) Polygon {
    // var edges = [];
    // var minX = (vertices.length > 0) ? vertices[0].X : undefined;
    // var minY = (vertices.length > 0) ? vertices[0].Y : undefined;
    // var maxX = minX;
    // var maxY = minY;

    // for (var i = 0; i < polygon.vertices.length; i++) {
    //     vertices[i].label = String(i);
    //     vertices[i].isReflex = isReflexVertex(polygon, i);
    //     var edge = {
    //         vertex1: vertices[i], 
    //         vertex2: vertices[(i + 1) % vertices.length], 
    //         polygon: polygon, 
    //         index: i
    //     };
    //     edge.outwardNormal = outwardEdgeNormal(edge);
    //     edge.inwardNormal = inwardEdgeNormal(edge);
    //     edges.push(edge);
    //     var x = vertices[i].X;
    //     var y = vertices[i].Y;
    //     minX = Math.min(x, minX);
    //     minY = Math.min(y, minY);
    //     maxX = Math.max(x, maxX);
    //     maxY = Math.max(y, maxY);
    // }                       
    
    // polygon.edges = edges;
    // polygon.minX = minX;
    // polygon.minY = minY;
    // polygon.maxX = maxX;
    // polygon.maxY = maxY;
    // polygon.closed = true;

    return Polygon{
		vertices: vertices,
	}
}

// based on http://local.wasp.uwa.edu.au/~pbourke/geometry/lineline2d/, edgeA => "line a", edgeB => "line b"

func edgesIntersection(edgeA Edge, edgeB Edge) domain.Position {
    den := (edgeB.vertex2.Y - edgeB.vertex1.Y) * (edgeA.vertex2.X - edgeA.vertex1.X) - (edgeB.vertex2.X - edgeB.vertex1.X) * (edgeA.vertex2.Y - edgeA.vertex1.Y)
    if den == 0 {
		return null  // lines are parallel or conincident
	}

    ua := ((edgeB.vertex2.X - edgeB.vertex1.X) * (edgeA.vertex1.Y - edgeB.vertex1.Y) - (edgeB.vertex2.Y - edgeB.vertex1.Y) * (edgeA.vertex1.X - edgeB.vertex1.X)) / den
    ub := ((edgeA.vertex2.X - edgeA.vertex1.X) * (edgeA.vertex1.Y - edgeB.vertex1.Y) - (edgeA.vertex2.Y - edgeA.vertex1.Y) * (edgeA.vertex1.X - edgeB.vertex1.X)) / den

    if ua < 0 || ub < 0 || ua > 1 || ub > 1 {
		return null
	}

    return domain.Position{
		X: edgeA.vertex1.X + ua * (edgeA.vertex2.X - edgeA.vertex1.X),
		Y: edgeA.vertex1.Y + ua * (edgeA.vertex2.Y - edgeA.vertex1.Y),
	}
}

func appendArc(vertices []domain.Position, center domain.Position, radius int, startVertex domain.Position, endVertex domain.Position, isPaddingBoundary bool) []domain.Position {
    const twoPI = math.twoPI
    startAngle := math.atan2(startVertex.Y - center.Y, startVertex.X - center.X)
    endAngle := math.atan2(endVertex.Y - center.Y, endVertex.X - center.X)
    if startAngle < 0 {
		startAngle += twoPI
	}
    if endAngle < 0 {
		endAngle += twoPI
	}
	arcSegmentCount := 5 // An odd number so that one arc vertex will be eactly arcRadius from center.
	angle := startAngle - endAngle
	if startAngle > endAngle {
		angle = startAngle + twoPI - endAngle
	}
	angle5 := -angle
	if isPaddingBoundary {
		angle5 = (twoPI - angle) / arcSegmentCount
	}

    vertices.push(startVertex);
    for (var i = 1; i < arcSegmentCount; ++i) {
        var angle = startAngle + angle5 * i;
        var vertex = {
            x: center.X + math.cos(angle) * radius,
            y: center.Y + math.sin(angle) * radius,
        };
        vertices.push(vertex);
    }
    vertices.push(endVertex);
}

function createOffsetEdge(edge, dx, dy)
{
    return {
        vertex1: {x: edge.vertex1.X + dx, y: edge.vertex1.Y + dy},
        vertex2: {x: edge.vertex2.X + dx, y: edge.vertex2.Y + dy}
    };
}

function createMarginPolygon(polygon)
{
    var offsetEdges = [];
    for (var i = 0; i < polygon.edges.length; i++) {
        var edge = polygon.edges[i];
        var dx = edge.outwardNormal.X * shapeMargin;
        var dy = edge.outwardNormal.Y * shapeMargin;
        offsetEdges.push(createOffsetEdge(edge, dx, dy));
    }

    var vertices = [];
    for (var i = 0; i < offsetEdges.length; i++) {
        var thisEdge = offsetEdges[i];
        var prevEdge = offsetEdges[(i + offsetEdges.length - 1) % offsetEdges.length];
        var vertex = edgesIntersection(prevEdge, thisEdge);
        if (vertex)
            vertices.push(vertex);
        else {
            var arcCenter = polygon.edges[i].vertex1;
            appendArc(vertices, arcCenter, shapeMargin, prevEdge.vertex2, thisEdge.vertex1, false);
        }
    }

    var marginPolygon = createPolygon(vertices);
    marginPolygon.offsetEdges = offsetEdges;
    return marginPolygon;
}

function createPaddingPolygon(polygon)
{
    var offsetEdges = [];
    for (var i = 0; i < polygon.edges.length; i++) {
        var edge = polygon.edges[i];
        var dx = edge.inwardNormal.X * shapePadding;
        var dy = edge.inwardNormal.Y * shapePadding;
        offsetEdges.push(createOffsetEdge(edge, dx, dy));
    }

    var vertices = [];
    for (var i = 0; i < offsetEdges.length; i++) {
        var thisEdge = offsetEdges[i];
        var prevEdge = offsetEdges[(i + offsetEdges.length - 1) % offsetEdges.length];
        var vertex = edgesIntersection(prevEdge, thisEdge);
        if (vertex)
            vertices.push(vertex);
        else {
            var arcCenter = polygon.edges[i].vertex1;
            appendArc(vertices, arcCenter, shapePadding, prevEdge.vertex2, thisEdge.vertex1, true);
        }
    }

    var paddingPolygon = createPolygon(vertices);
    paddingPolygon.offsetEdges = offsetEdges;
    return paddingPolygon;
}

function computeAll()
{
    polygon = createPolygon(polygon.vertices);
    marginPolygon = createMarginPolygon(polygon);
    paddingPolygon = createPaddingPolygon(polygon);
}

function init() 
{
    var polygonVertices =  [{x: 143, y: 327}, {x: 80, y: 236}, {x: 151, y: 148}, {x: 454, y: 69}, {x: 560, y: 320}];
    polygon = createPolygon(polygonVertices);


    computeAll();
}