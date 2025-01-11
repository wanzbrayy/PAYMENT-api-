// Load Pie Chart using Chart.js
const ctx = document.getElementById("pieChart").getContext("2d");
const pieChart = new Chart(ctx, {
  type: "pie",
  data: {
    labels: ["Aktif", "Nonaktif", "Tertunda"],
    datasets: [
      {
        label: "Status Pengguna",
        data: [60, 30, 10], // Data awal bisa diganti dengan API call
        backgroundColor: ["#2ecc71", "#e74c3c", "#f39c12"],
        hoverOffset: 4,
      },
    ],
  },
  options: {
    responsive: true,
    plugins: {
      legend: {
        position: "top",
      },
    },
  },
});

// Fetch Data from API
async function fetchDashboardData() {
  try {
    const response = await fetch("http://localhost:8080/financial-status"); // Ganti dengan URL API yang sesuai
    const data = await response.json();

    // Update Statistik
    document.getElementById("total-income").innerText = `Rp${data.total_income.toLocaleString()}`;
    document.getElementById("total-expenses").innerText = `Rp${data.total_expenditure.toLocaleString()}`;
    document.getElementById("total-balance").innerText = `Rp${(data.total_income - data.total_expenditure).toLocaleString()}`;

    // Update Pie Chart Data (contoh data pengguna)
    pieChart.data.datasets[0].data = [50, 40, 10]; // Update dengan data yang relevan
    pieChart.update();

    // Update Notifikasi dan tabel pengguna
    const userTableBody = document.querySelector("#user-table tbody");
    userTableBody.innerHTML = ""; // Clear existing rows
    const users = data.users || []; // Ambil data pengguna jika ada

    users.forEach(user => {
      const row = document.createElement("tr");
      row.innerHTML = `
        <td>${user.username}</td>
        <td>${user.status}</td>
        <td><button class="btn-edit" data-id="${user.id}">Edit</button></td>
      `;
      userTableBody.appendChild(row);

      // Edit button functionality
      row.querySelector(".btn-edit").addEventListener("click", () => {
        openEditProfile(user);
      });
    });

  } catch (error) {
    console.error("Error fetching dashboard data:", error);
  }
}

// Fetch data when page loads
document.addEventListener("DOMContentLoaded", fetchDashboardData);

// Toggle Profile Modal
document.getElementById("profile-btn").onclick = () => {
  document.getElementById("profile-modal").style.display = "flex";
};

document.getElementById("close-profile").onclick = () => {
  document.getElementById("profile-modal").style.display = "none";
};

// Logout Button
document.getElementById("logout-btn").onclick = () => {
  alert("Anda berhasil logout");
  window.location.href = "index.html"; // Redirect to login page
};

// Open Edit Profile Modal with User Data
function openEditProfile(user) {
  document.getElementById("profile-name").value = user.username;
  document.getElementById("profile-email").value = user.email;
  document.getElementById("profile-modal").style.display = "flex";
}

// Save edited profile
document.getElementById("profile-form").onsubmit = async (e) => {
  e.preventDefault();
  const username = document.getElementById("profile-name").value;
  const email = document.getElementById("profile-email").value;

  try {
    const response = await fetch("http://localhost:8080/update-profile", { // Ganti URL dengan API endpoint
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, email }),
    });

    const result = await response.json();

    if (result.success) {
      alert("Profil berhasil diperbarui!");
      document.getElementById("profile-modal").style.display = "none";
      fetchDashboardData(); // Refresh data
    } else {
      alert("Terjadi kesalahan saat memperbarui profil.");
    }
  } catch (error) {
    console.error("Error updating profile:", error);
  }
};
