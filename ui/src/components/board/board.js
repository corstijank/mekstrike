const a = (2 * Math.PI) / 6;

export const hexSize = 60; // Size of one edge of the hexagon
export const hexWidth = hexSize * (1 + Math.cos(a));
export const hexHeight = hexSize * Math.sin(a) + hexSize * Math.sin(a);

export function colToCenterX(col){
    return col * hexWidth - hexSize * 0.5;
}

export function rowToCenterY(row,col){
	var y = row * hexHeight - hexSize * Math.sin(a);
	if (col % 2 == 0) {
		y += 1 ** (row + 1) * hexSize * Math.sin(a);
	}
	return  y ;
}

export function getHexCenter(row, col) {
	var x = col * hexWidth - hexSize * 0.5;
	var y = row * hexHeight - hexSize * Math.sin(a);
	if (col % 2 == 0) {
		y += 1 ** (row + 1) * hexSize * Math.sin(a);
	}

	return { x, y };
}
