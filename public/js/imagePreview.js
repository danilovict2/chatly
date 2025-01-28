const enableImagePreview = (e) => {
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.onloadend = () => {
        const previewContainer = document.createElement("div");
        previewContainer.classList.add("preview-container", "mb-3", "flex", "items-center", "gap-2");

        previewContainer.innerHTML = `
            <div class="relative">
                <img
                  src="${reader.result}"
                  alt="Preview"
                  class="w-20 h-20 object-cover rounded-lg border border-zinc-700"
                />
                <button
                  onclick="disableImagePreview(this)"
                  class="absolute -top-1.5 -right-1.5 w-5 h-5 rounded-full bg-base-300
                  flex items-center justify-center"
                  type="button"
                >
                  <i class="fa-solid fa-x size-3"></i>
                </button>
            </div>
        `;

        const inputContainer = document.getElementById("input-container");
        inputContainer.prepend(previewContainer);
    };

    reader.readAsDataURL(file);
}

const disableImagePreview = (button) => {
    const previewContainer = button.closest('.preview-container');
    previewContainer.remove()

    const fileInput = document.getElementById('message-image-input');
    if (fileInput) fileInput.value = '';
}

const clickImageInput = () => {
    document.getElementById("message-image-input").click()
}

document.getElementById("message-image-input").addEventListener("change", enableImagePreview);