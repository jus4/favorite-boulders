document.addEventListener("DOMContentLoaded", function () {
  const mapContainer = document.getElementById('map');

  if (!mapContainer) return;

  const observer = new IntersectionObserver(async (entries, observer) => {
    if (entries[0].isIntersecting) {
      observer.disconnect();

      const module = await import('./modules/terrainMap.ts');
      const newMap = new module.default();
      newMap.init();
    }
  });

  observer.observe(mapContainer);
});
