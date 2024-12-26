const submitButton = document.getElementById('submitButton');
const noteInput = document.getElementById('noteInput');
const noteOutput = document.getElementById('noteOutput');
const noteIdSpan = document.getElementById('noteId');

const fetchButton = document.getElementById('fetchButton');
const noteIdInput = document.getElementById('noteIdInput');

// Отправка новой заметки
submitButton.addEventListener('click', async () => {
    const noteText = noteInput.value.trim();

    if (noteText === '') {
        alert('Please enter a note!');
        return;
    }

    // Получаем токен капчи
    const captchaToken = document.querySelector('input[name="captcha_token"]')?.value;


    try {
        //const response = await fetch('http://localhost:8080/notes', {
        const response = await fetch('http://51.250.47.117:8080/notes', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                text: noteText,
                captcha_token: captchaToken  // Отправляем токен капчи
            })
        });

        if (!response.ok) {
            throw new Error('Failed to save note');
        }

        const result = await response.json();
        console.info(result);

        // Отображаем данные
        noteIdInput.value  = result.Id;

        noteInput.value = ''; // Очищаем поле ввода
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to save note. Please try again later.');
    }
});

// Получение заметки по ID
fetchButton.addEventListener('click', async () => {
    const noteId = noteIdInput.value.trim();

    if (noteId === '') {
        alert('Please enter a Note ID!');
        return;
    }

    try {
        //const response = await fetch(`http://localhost:8080/notes/${noteId}`, {
        const response = await fetch(`http://51.250.47.117:8080/notes/${noteId}`, {
            method: 'GET',
        });

        if (!response.ok) {
            throw new Error('Note not found');
        }

        const result = await response.json();

        // Заполняем поле ввода заметки текстом из ответа
        noteInput.value = result.text;
        
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to fetch note. Please check the Note ID and try again.');
    }
});


// Инициализация капчи
smartCaptcha.render('captcha-container', {
    sitekey: 'ysc1_GJTNDhHixFgjbWbP2ANiBudNjxf7uuEcGqninafg1d089560', // Ваш sitekey
    callback: function (token) {
        // Капча решена, добавляем токен в форму
        var tokenInput = document.createElement('input');
        tokenInput.type = 'hidden';
        tokenInput.name = 'captcha_token';
        tokenInput.value = token;
        document.getElementById('captcha-container').appendChild(tokenInput);
    }
});
