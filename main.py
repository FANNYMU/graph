import os
import datetime
import time
import random
import subprocess
from colorama import init, Fore, Style

# Inisialisasi colorama
init()

def buat_fake_commit(target_commits=100):
    print(f"{Fore.CYAN}=== GITHUB FAKE COMMIT GENERATOR ==={Style.RESET_ALL}")
    print(f"{Fore.YELLOW}Target commit: {target_commits}{Style.RESET_ALL}")
    
    # Pastikan file test.txt ada
    if not os.path.exists("test.txt"):
        with open("test.txt", "w") as f:
            f.write(f"Commit for {datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
    
    commits_selesai = 0
    
    try:
        # Mulai loop commit
        while commits_selesai < target_commits:
            # Buat waktu commit dengan format yang benar
            waktu_sekarang = datetime.datetime.now()
            format_waktu = waktu_sekarang.strftime("%Y-%m-%d %H:%M:%S")
            
            # Tambahkan konten ke file
            with open("test.txt", "a") as f:
                f.write(f"Commit for {datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
            
            # Jalankan perintah git
            subprocess.run(["git", "add", "test.txt"], check=True)
            subprocess.run(["git", "commit", "-m", f"Commit {commits_selesai + 1} - {format_waktu}"], check=True)
            
            # Update counter dan tampilkan progress
            commits_selesai += 1
            
            # Tampilan progress bar
            progress = int((commits_selesai / target_commits) * 30)
            print(f"\r{Fore.GREEN}[{'#' * progress}{' ' * (30-progress)}] {commits_selesai}/{target_commits} commits {Fore.CYAN}({int(commits_selesai/target_commits*100)}%){Style.RESET_ALL}", end="")
            
            # Tambahkan jeda random agar terlihat lebih natural
            # time.sleep(random.uniform(0.5, 2.0))
        
        print("\n")
        print(f"{Fore.GREEN}✓ {Style.BRIGHT}Selesai! {commits_selesai} commits berhasil dibuat.{Style.RESET_ALL}")
        print(f"{Fore.YELLOW}Jangan lupa untuk melakukan git push untuk mengirim commit ke repository.{Style.RESET_ALL}")
        print(f"{Fore.CYAN}git push origin <nama-branch>{Style.RESET_ALL}")
        
    except KeyboardInterrupt:
        print(f"\n\n{Fore.RED}Program dihentikan oleh pengguna.{Style.RESET_ALL}")
        print(f"{Fore.YELLOW}{commits_selesai} commits telah dibuat sebelum dihentikan.{Style.RESET_ALL}")
    except Exception as e:
        print(f"\n\n{Fore.RED}Error: {str(e)}{Style.RESET_ALL}")

if __name__ == "__main__":
    try:
        jumlah_commit = int(input(f"{Fore.CYAN}Masukkan jumlah commit yang diinginkan: {Style.RESET_ALL}"))
        buat_fake_commit(jumlah_commit)
    except ValueError:
        print(f"{Fore.RED}Error: Masukkan angka yang valid!{Style.RESET_ALL}")
