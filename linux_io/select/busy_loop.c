#include <stdio.h>
#include <stdlib.h>
#include <sys/select.h>
#include <unistd.h>

int main() {
    fd_set read_fds;
    FD_ZERO(&read_fds);
    FD_SET(STDIN_FILENO, &read_fds);  // Monitor stdin

    while (1) {
        // Wait for data to be available on stdin
        int ret = select(STDIN_FILENO + 1, &read_fds, NULL, NULL, NULL);
        if (ret == -1) {
            perror("select");
            exit(EXIT_FAILURE);
        }

        // Check if stdin is readable
        if (FD_ISSET(STDIN_FILENO, &read_fds)) {
            printf("Data is available on stdin, but not reading it.\n");

            // Uncomment the following code to read the data and avoid busy loop
            /*
            char buf[1024];
            ssize_t n = read(STDIN_FILENO, buf, sizeof(buf));
            if (n > 0) {
                printf("Received data: %.*s", (int)n, buf);
            } else {
                break;  // End of input
            }
            */
        }
    }

    return 0;
}