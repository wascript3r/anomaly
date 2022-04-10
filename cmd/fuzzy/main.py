import numpy as np
import skfuzzy as fuzz
import matplotlib.pyplot as plt

x_day_time = np.arange(0, 24, 1)
x_week_day = np.arange(1, 8, 1)
x_last_imsi_calls = np.arange(0, 41, 1)
x_last_msc_calls = np.arange(0, 121, 1)
x_prob = np.arange(0, 101, 1)

day_time_lo = fuzz.trapmf(x_day_time, [0, 0, 5, 8])
day_time_md = fuzz.trapmf(x_day_time, [6, 9, 11, 14])
day_time_hi = fuzz.trapmf(x_day_time, [11, 14, 17, 19])
day_time_vh = fuzz.trapmf(x_day_time, [17, 19, 22, 23])
day_time_eh = fuzz.trapmf(x_day_time, [21, 23, 23, 23])

week_day_lo = fuzz.trapmf(x_week_day, [1, 1, 4, 6])
week_day_md = fuzz.trapmf(x_week_day, [4, 6, 7, 7])

last_imsi_calls_lo = fuzz.trapmf(x_last_imsi_calls, [0, 0, 5, 8])
last_imsi_calls_md = fuzz.trapmf(x_last_imsi_calls, [5, 9, 15, 20])
last_imsi_calls_hi = fuzz.trapmf(x_last_imsi_calls, [14, 21, 40, 40])

last_msc_calls_lo = fuzz.trapmf(x_last_msc_calls, [0, 0, 30, 50])
last_msc_calls_md = fuzz.trapmf(x_last_msc_calls, [30, 55, 80, 100])
last_msc_calls_hi = fuzz.trapmf(x_last_msc_calls, [82, 105, 120, 120])

prob_lo = fuzz.trapmf(x_prob, [0, 0, 20, 38])
prob_md = fuzz.trapmf(x_prob, [25, 40, 50, 60])
prob_hi = fuzz.trapmf(x_prob, [50, 60, 70, 80])
prob_vh = fuzz.trapmf(x_prob, [72, 85, 100, 100])

fig1, (ax1, ax2) = plt.subplots(nrows=2, ncols=2, figsize=(9, 12))
fig2, (ax3, ax4) = plt.subplots(nrows=2, ncols=2, figsize=(9, 12))

ax1[0].plot(x_day_time, day_time_lo, "blue", linewidth=2, label="Naktis")
ax1[0].plot(x_day_time, day_time_md, "purple", linewidth=2, label="Rytas")
ax1[0].plot(x_day_time, day_time_hi, "magenta", linewidth=2, label="Diena")
ax1[0].plot(x_day_time, day_time_vh, "green", linewidth=2, label="Vakaras")
ax1[0].plot(x_day_time, day_time_eh, "blue", linewidth=2, label="Naktis")
ax1[0].set_title("Paros laikas")
ax1[0].xaxis.set_ticks(np.arange(min(x_day_time), max(x_day_time) + 1, 1.0))
ax1[0].legend()

ax1[1].plot(x_week_day, week_day_lo, "blue", linewidth=2, label="Darbo diena")
ax1[1].plot(x_week_day, week_day_md, "purple", linewidth=2, label="Savaitgalis")
ax1[1].set_title("Savaitės diena")
ax1[1].xaxis.set_ticks(np.arange(min(x_week_day), max(x_week_day) + 1, 1.0))
ax1[1].legend()

ax2[0].plot(x_last_imsi_calls, last_imsi_calls_lo, "blue", linewidth=2, label="Mažas")
ax2[0].plot(x_last_imsi_calls, last_imsi_calls_md, "purple", linewidth=2, label="Vidutinis")
ax2[0].plot(x_last_imsi_calls, last_imsi_calls_hi, "magenta", linewidth=2, label="Didelis")
ax2[0].set_title("Paskutinės valandos skambučių kiekis iš IMSI")
ax2[0].xaxis.set_ticks(np.arange(min(x_last_imsi_calls), max(x_last_imsi_calls) + 1, 2.0))
ax2[0].legend()

ax2[1].plot(x_last_msc_calls, last_msc_calls_lo, "blue", linewidth=2, label="Mažas")
ax2[1].plot(x_last_msc_calls, last_msc_calls_md, "purple", linewidth=2, label="Vidutinis")
ax2[1].plot(x_last_msc_calls, last_msc_calls_hi, "magenta", linewidth=2, label="Didelis")
ax2[1].set_title("Paskutinės valandos skambučių kiekis į MSC")
ax2[1].xaxis.set_ticks(np.arange(min(x_last_msc_calls), max(x_last_msc_calls) + 1, 5.0))
ax2[1].legend()

ax3[0].plot(x_prob, prob_lo, "blue", linewidth=2, label="Maža")
ax3[0].plot(x_prob, prob_md, "purple", linewidth=2, label="Vidutinė")
ax3[0].plot(x_prob, prob_hi, "magenta", linewidth=2, label="Didelė")
ax3[0].plot(x_prob, prob_vh, "red", linewidth=2, label="Labai didelė")
ax3[0].set_title("DoS atakos tikimybė")
ax3[0].xaxis.set_ticks(np.arange(min(x_prob), max(x_prob) + 1, 5.0))
ax3[0].legend()

plt.show()
