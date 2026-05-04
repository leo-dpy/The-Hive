(function () {
  function updateCounter(counter) {
    const id = counter.getAttribute("data-for");
    const target = document.getElementById(id);
    if (!target) return;
    const max = target.getAttribute("maxlength") || "∞";
    counter.textContent = `${target.value.length} / ${max}`;
  }

  document.querySelectorAll(".counter[data-for]").forEach((counter) => {
    const id = counter.getAttribute("data-for");
    const target = document.getElementById(id);
    if (!target) return;
    updateCounter(counter);
    target.addEventListener("input", () => updateCounter(counter));
  });
})();
