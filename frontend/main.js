function setView(element, data) {
  const list = document.createElement("ul")
  element.appendChild(list)
  data.teams.forEach(team => {
    const item = document.createElement("li")
    item.innerText = `${team.name} (${team.short_name})`
    list.appendChild(item)
  })
}

export {
  setView
}