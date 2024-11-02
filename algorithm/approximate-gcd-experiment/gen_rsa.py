"""
This module provides functions to generate RSA key pairs and sign messages.
The functions are similar to the ones that are used in Sign (ISITDTU 2024 QUALS CTF).
"""
import os
from Crypto.Util.number import getPrime, GCD
from Crypto.Signature import PKCS1_v1_5
from Crypto.PublicKey import RSA
from Crypto.Hash import SHA256


def genkey(e=11):
    """
    Generates an RSA key pair with the public exponent e.
    """
    while True:
        p = getPrime(1024)
        q = getPrime(1024)
        if GCD(p-1, e) == 1 and GCD(q-1, e) == 1:
            break
    n = p*q
    d = pow(e, -1, (p-1)*(q-1))
    return RSA.construct((int(n), int(e), int(d)))


def gensig(key: RSA.RsaKey) -> bytes:
    """
    Generates a PKCS#1 v1.5 RSA signature for a random message using the given RSA key.
    """
    m = os.urandom(256)
    h = SHA256.new(m)
    s = PKCS1_v1_5.new(key).sign(h)
    return s
