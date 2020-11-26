import argparse
import subprocess
import sys

def parser():
    p = argparse.ArgumentParser()
    p.add_argument(
        "--host",
        required=True,
    )
    p.add_argument(
        "--dir", "-d",
        default="/",
    )

    return p


def main(argv=[]):
    args = parser().parse_args(argv)

    host = args.host

    ssh_ls = subprocess.Popen(["ssh", host, "ls {}".format(args.dir)], stdout=subprocess.PIPE)
    shuf = subprocess.Popen(["shuf"], stdin=ssh_ls.stdout, stdout=subprocess.PIPE)
    head1 = subprocess.Popen(["head", "-n", "1"], stdin=shuf.stdout, stdout=subprocess.PIPE)

    ssh_ls.stdout.close()
    shuf.stdout.close()
    output, err = head1.communicate()
    ssh_ls.wait()
    shuf.wait()

    if err:
        print(err)
        return -1

    if ssh_ls.returncode != 0:
        return ssh_ls.returncode

    if shuf.returncode != 0:
        return shuf.returncode

    if head1.returncode != 0:
        return head1.returncode

    print(output.rstrip())
    return 0


if __name__ == '__main__':
    sys.exit(main(sys.argv[1:]))
    