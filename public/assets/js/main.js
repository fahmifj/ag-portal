$(document).ready(() => {
	loadVMs().then(controlVM);
});


const loadVMs = async () => {
	await fetch(`/vm/fetch`, {
		method: 'GET',
	})
		.then((response) => response.json()
			.then((resp) => {
				if (response.status == 200) {
					$(resp).each((k, entry) => {
						let ipHTML, vmNameHTML, statusHTML;
						vmNameHTML = `<div id="${entry.name}" class="vm-name text" >${entry.name}</div>`
						statusHTML = `<div><span>Status&nbsp;</span><span id="vm-status" class="status" >${entry.status}</span></div>`
						ipHTML = ""

						if (entry.publicIp !== undefined) {
							ipHTML = `<div><span>Public IP&nbsp;</span><span class="ip" >${entry.publicIp}</span></div>`;
						}

						entryHTML = `<div class="card-item">` + vmNameHTML +
							`<div class="info">` + statusHTML + ipHTML + `</div>` +
							`<a class="btn-control" href="javascript:void(0);" >Start</a>`;

						$('#entries-container').append(entryHTML);
					})
				}

			}))
		.catch((error) => {

		});
}

const controlVM = () => {

	console.log("VM list has been loaded");

	$('.card-item').each((i, ci) => {
		
		button = ci.lastElementChild;
		button.addEventListener("click", (e) => {
			e.preventDefault();
			
			const vmName = ci.childNodes[0].textContent;
			const btn = ci.childNodes[2];
			const action = btn.textContent.toLowerCase();
			fetch(`/vm/${vmName}?action=${action}`, {
				method: 'POST',				
			});
			
			switch (action) {
				case "start":
					btn.textContent = "Stop";
					btn.style.backgroundColor = 'var(--red)';
					break;
				case "stop":
					btn.textContent = "Start";
					btn.style.backgroundColor = 'var(--green)';
					break;
			}

		});

	});

}

