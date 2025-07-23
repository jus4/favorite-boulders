const proj4jConfig = () => {
  const proj4 = globalThis.proj4
  const ol = globalThis.ol

  proj4.defs("EPSG:3067", "+proj=utm +zone=35 +ellps=GRS80 +units=m +no_defs");
  ol.proj.proj4.register(proj4);

  const projection = ol.proj.get('EPSG:3067');
  const extent = [-548576.000000, 6291456.000000, 1548576.000000, 8388608.000000];

  const resolutions = [
    8192, 4096, 2048, 1024, 512, 256, 128, 64, 32, 16, 8, 4, 2, 1, 0.5
  ];
  const matrixIds = resolutions.map((_, i) => i);
  return {
    projection,
    extent,
    resolutions,
    matrixIds
  }
}

export default proj4jConfig()
